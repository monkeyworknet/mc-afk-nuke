package main

import (
	"encoding/json"
	"fmt"
	"github.com/monkeyworknet/mc-afk-nuke/config"
	"github.com/nanobox-io/golang-scribble"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Playersdb struct {
	// Create structure for JSON DB files.

	Name      string `json:"name"`
	Filename  string `json:"filename"`
	Totalmove int    `json:"totalmove"`
	Afkindex  int    `json:"afkindex"`
}

var (
	logfilename = "./mc-afk-nuke.log"
	logprefix   = "MCAFK : "
	out, _      = os.OpenFile(logfilename, os.O_APPEND|os.O_WRONLY, 0644)
	flag        = log.LstdFlags | log.Lshortfile
	newLog      = log.New(out, logprefix, flag)
)

func main() {

	// Read Config
	err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error Loading Config Functions")
		os.Exit(2)

	}

	// Setup Logging

	_, err = os.Stat(logfilename)
	if os.IsNotExist(err) {
		file, _ := os.Create(logfilename)
		defer file.Close()
	}

	// Grab active users on server

	resp, err := http.Get("https://api.minetools.eu/query/" + config.Servername + "/" + config.Serverport)

	if err != nil {
		newLog.Println("FATAL Error getting Playerlist (minetools down?)  -  ", err)
		os.Exit(2)
	}

	defer resp.Body.Close()

	// Check to see if anyone is online before we do anything

	var NumberPlayers struct {
		Playercount int `json:"Players"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&NumberPlayers); err != nil {
		newLog.Println("Error Decoding PlayerCountlist  -  ", err)
		os.Exit(4)
	}

	if NumberPlayers.Playercount == 0 {
		fmt.Println("No One Online")
		os.Exit(0)
	}

	// Grab the player list

	resp, err = http.Get("https://api.minetools.eu/query/" + config.Servername + "/" + config.Serverport)

	if err != nil {
		newLog.Println("FATAL Error getting Playerlist (minetools down?)  -  ", err)
		os.Exit(2)
	}

	defer resp.Body.Close()

	var ActivePlayers struct {
		Playerlist  []string `json:"Playerlist"`
		Playercount int      `json:"Players"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ActivePlayers); err != nil {
		newLog.Println("Error Decoding Playerlist  -  ", err)
		os.Exit(4)
	}

	afkkickvalue := randInt(config.Afkkickvaluemin, config.Afkkickvaluemax)
	fmt.Println("INFO - AFK Kick Value Set to ", afkkickvalue)
	fmt.Println("PLAYER NAME | UUID | CURRENT MOVE | HISTORICAL MOVE | AFK COUNTER")

	// if player list is greater than 1 then itterate over playerlist

	if ActivePlayers.Playercount > 0 {
		for _, player := range ActivePlayers.Playerlist {
			uuidfilename, uuid := getUUID(player)
			fullpath := config.Path + uuidfilename
			currenttotalmovement := playerMovementStats(fullpath)
			newtotmov, newafk := jSONdbwork(player, uuid, currenttotalmovement)
			fmt.Println(player, uuid, currenttotalmovement, newtotmov, newafk)
			if afkkickvalue <= newafk {
				// run kick function
				_ = kickPlayer(player, uuid)
			}
		}

	} else {
		fmt.Println("Server Empty")
	}

}

func randInt(min int, max int) int {

	//  Generate the random AFK kick level

	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}

func getUUID(name string) (string, string) {

	// give string [name] return string [uuid] in the form of a .json file on system and normal UUID
	// grab trimmed UUID from Mojang for playername

	db, err := scribble.New(config.Database, nil)
	if err != nil {
		newLog.Println("FATAL Error creating db", err)
		os.Exit(3)
	}
	playerread := Playersdb{}
	ut := "0"
	if err := db.Read(config.Collection, name, &playerread); err != nil {
		newLog.Println("INFO - couldn't UUID for "+name+" in DB Looking Online", err)
		fmt.Println("Looking UUID Up online..")

		var UUID struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}

		resp, err := http.Get("https://api.minetools.eu/uuid/" + name)
		if err != nil {
			newLog.Println("Error getting UUID (minetools api?)  -  ", err)
			return "", ""
		}

		if err := json.NewDecoder(resp.Body).Decode(&UUID); err != nil {
			newLog.Println("Error Decoding UUID  -  ", err)
			return "", ""
		}
		ut = UUID.Id
	}
	if ut == "0" {
		ut = playerread.Filename
	}

	// insert hypens as per file name and add extension

	index1 := 8
	index2 := 12
	index3 := 16
	index4 := 20
	dash := "-"
	ext := ".json"

	filename := ut[:index1] + dash + ut[index1:index2] + dash + ut[index2:index3] + dash + ut[index3:index4] + dash + ut[index4:] + ext

	return filename, ut
}

func playerMovementStats(file string) int {

	// Take in filename return init representing total movemement on server stats

	raw, err := ioutil.ReadFile(file)

	if err != nil {
		newLog.Println("ERROR - can't read file ", file)
		return 0
	}

	var stats struct {
		Walked   int `json:"stat.walkOneCm"`
		Crouched int `json:"stat.crouchedOneCm"`
		Sprint   int `json:"stat.sprintOneCm"`
		Swim     int `json:"stat.swimOneCm"`
		Fall     int `json:"stat.fallOneCm"`
		Climb    int `json:"stat.climbOneCm"`
		Fly      int `json:"stat.flyOneCm"`
		Dive     int `json:"stat.diveOneCm"`
		Minecart int `json:"stat.minecartOneCm"`
		Boat     int `json:"stat.boatOneCm"`
		Pig      int `json:"stat.pigOneCm"`
		Horse    int `json:"stat.horseOneCm"`
		Aviate   int `json:"stat.aviateOneCm"`
	}

	if err := json.Unmarshal(raw, &stats); err != nil {
		newLog.Println("Error - Can't understand stats json file ", file)
		return 0
	}

	TotalMovementReturn := stats.Walked + stats.Crouched + stats.Sprint + stats.Swim + stats.Fall + stats.Climb + stats.Fly + stats.Dive + stats.Minecart + stats.Boat + stats.Pig + stats.Horse + stats.Aviate

	return TotalMovementReturn
}

func jSONdbwork(name string, filename string, totalmove int) (int, int) {
	// read in name, totalmove, and uuidfilename to use incase entry doesn't exist in database
	// if entry doesn't exist insert it with an afkindex of 0
	// if entry DOES exist compare movement entries, if they are the same increase afkindex, if different reset to 0
	// return read in total move int and afk index int

	db, err := scribble.New(config.Database, nil)
	if err != nil {
		newLog.Println("FATAL Error creating db", err)
		os.Exit(3)
	}

	afkvalue := 0
	playerread := Playersdb{}
	playerwrite := Playersdb{Name: name, Filename: filename, Totalmove: totalmove, Afkindex: afkvalue}
	if err := db.Read(config.Collection, name, &playerread); err != nil {
		newLog.Println("INFO - couldn't find history in DB ", name, err)
		if err := db.Write(config.Collection, name, playerwrite); err != nil {
			newLog.Println("ERROR - couldn't write initial entry to DB", name, err)
		}

	}

	if playerread.Totalmove == totalmove {
		afkvalue = playerread.Afkindex + 1
		playerwrite = Playersdb{Name: name, Filename: filename, Totalmove: totalmove, Afkindex: afkvalue}
		newLog.Println("INFO - Incremented AFK value for ", name)
	} else {
		afkvalue = 0
		playerwrite = Playersdb{Name: name, Filename: filename, Totalmove: totalmove, Afkindex: afkvalue}
		newLog.Println("INFO - AFK Values Reset for ", name)
	}

	if err := db.Write(config.Collection, name, playerwrite); err != nil {
		newLog.Println("Error -  couldn't write update to DB ", err)
	}

	return playerread.Totalmove, afkvalue
}

func kickPlayer(name string, filename string) error {
	// OS Exec the actual kick command and reset kickindex for user in DB.
	//  If you use a different init script than what is listed in the README you will need to change this.

	db, err := scribble.New(config.Database, nil)
	if err != nil {
		fmt.Println("Error creating db", err)
	}

	playerwrite := Playersdb{Name: name, Filename: filename, Totalmove: 0, Afkindex: 0}
	fmt.Println("Kicking Player for exceeding AFK timer ", name)
	kickcmd := "/usr/sbin/service"
	kickargs := []string{"minecraft", "command kick", name, config.Kickreason}
	newLog.Println("Kicking Player ", name)
	if err := exec.Command(kickcmd, kickargs...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		newLog.Print("Something went wrong when kicking", name, err)
	}

	if err := db.Write(config.Collection, name, playerwrite); err != nil {
		newLog.Println("Error - couldn't write update to DB ", err)
	}
	return nil
}
