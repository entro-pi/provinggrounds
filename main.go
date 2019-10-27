package main

import (
	"bufio"
	"os"
	"context"
	"time"
	"fmt"
	"strconv"
	"strings"
	"math/rand"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/SolarLune/dngn"
)
func clear() {
	for i := 0;i < 50;i++ {
		fmt.Println("                                                              ")
	}
}
type Descriptions struct {
	BATTLESPAM int
	ROOMDESC int
	PLAYERDESC int
	ROOMTITLE int
}
type Chat struct {
	User Player
	Message string
	Time time.Time
}
type Space struct{
	Room dngn.Room
	Vnums string
	Zone string
	ZonePos []int
	ZoneMap [][]int
	Vnum int
	Desc string
	Mobiles []int
	Items []int
	CoreBoard string
	Exits Exit
	Altered bool
}
type Exit struct {
	North int
	South int
	East int
	West int
	NorthWest int
	NorthEast int
	SouthWest int
	SouthEast int
	Up int
	Down int
}
func initDigRoom(digFrame [][]int, zoneVnums string, zoneName string, play Player, vnum int) (Space, int) {
	var dg Space
	dg.Vnums = zoneVnums
	dg.Zone = zoneName
	dg.ZonePos = make([]int, 2)
	dg.ZoneMap = digFrame
	//todo directions
	vnum += 1
	dg.Vnum = vnum
	dg.Altered = true
	dg.Desc = "Nothing but some cosmic rays"
	for len(strings.Split(dg.Desc, "\n")) < 8 {
		dg.Desc += "\n"
	}
	return dg, vnum
}

type Player struct {
	Name string
	Title string
	Inventory []int
	Equipment []int
	CoreBoard string
	CurrentRoom Space

	Str int
	Int int
	Dex int
	Wis int
	Con int
	Cha int
}


func DescribePlayer(play Player) {
	fmt.Printf("\033[40;0H")
	fmt.Println("======================")
	fmt.Println("\033[38:2:0:200:0mStrength     :\033[0m", play.Str)
	fmt.Println("\033[38:2:0:200:0mIntelligence :\033[0m", play.Int)
	fmt.Println("\033[38:2:0:200:0mDexterity    :\033[0m", play.Dex)
	fmt.Println("\033[38:2:0:200:0mWisdom       :\033[0m", play.Wis)
	fmt.Println("\033[38:2:0:200:0mConstitution :\033[0m", play.Con)
	fmt.Println("\033[38:2:0:200:0mCharisma     :\033[0m", play.Cha)
	fmt.Println("======================")
}
func updateRoom(play Player, populated []Space) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"vnum": bson.M{"$eq":play.CurrentRoom.Vnum}}
	collection := client.Database("zones").Collection("Spaces")
	update := bson.M{"$set": bson.M{"vnums":populated[play.CurrentRoom.Vnum].Vnums,
		"zone":populated[play.CurrentRoom.Vnum].Zone,"vnum":populated[play.CurrentRoom.Vnum].Vnum,
		 "desc":populated[play.CurrentRoom.Vnum].Desc,"exits": populated[play.CurrentRoom.Vnum].Exits,
			"mobiles": populated[play.CurrentRoom.Vnum].Mobiles, "items": populated[play.CurrentRoom.Vnum].Items,
			 "altered": true,"zonepos":populated[play.CurrentRoom.Vnum].ZonePos, "zonemap": populated[play.CurrentRoom.Vnum].ZoneMap }}

	result, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		panic(err)
	}
	fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
}

func updateZoneMap(play Player, populated []Space) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	filter := bson.M{"zone": bson.M{"$eq":play.CurrentRoom.Zone}}
	collection := client.Database("zones").Collection("Spaces")
	findOptions := options.Find()
	result, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		panic(err)
	}
	defer result.Close(context.Background())
	for result.Next(context.Background()) {
		var current Space
		err := result.Decode(&current)
		if err != nil {
			panic(err)
		}
		filter = bson.M{"zone": bson.M{"$eq":current.Zone}}
		update := bson.M{"$set": bson.M{"zonepos":populated[current.Vnum].ZonePos, "zonemap": populated[current.Vnum].ZoneMap }}

		result, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			panic(err)
		}
		fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
	}

}



func InitPlayer(name string) Player {
	var play Player
	var inv []int
	var equ []int
	inv = append(inv, 1)
	equ = append(equ, 1)
	play.Name = name
	play.Title = "The Unknown"
	play.Inventory = inv
	play.Equipment = equ
	play.Str = 1
	play.Int = 1
	play.Dex = 1
	play.Wis = 1
	play.Con = 1
	play.Cha = 1
	return play
}

func InitZoneSpaces(SpaceRange string, zoneName string, desc string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("zones").Collection("Spaces")
	vnums := strings.Split(SpaceRange, "-")
	vnumStart, err := strconv.Atoi(vnums[0])
	if err != nil {
		panic(err)
	}

	vnumEnd, err := strconv.Atoi(vnums[1])
	if err != nil {
		panic(err)
	}
	for i := vnumStart;i < vnumEnd;i++ {
		var mobiles []int
		var items []int
		mobiles = append(mobiles, 0)
		items = append(items, 0)
		_, err = collection.InsertOne(context.Background(), bson.M{"vnums":SpaceRange,"zone":zoneName,"vnum":i, "desc":desc,
							"mobiles": mobiles, "items": items })
	}
	if err != nil {
		panic(err)
	}
}
func PopulateAreaBuild() []Space {
	areas := make([]Space, 150)
	return areas
}

func PopulateAreas() []Space {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	var Spaces []Space
	collection := client.Database("zones").Collection("Spaces")
	results, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	for results.Next(context.Background()) {

			var Space Space
			err := results.Decode(&Space)
			if err != nil {
				panic(err)
			}
			Spaces = append(Spaces, Space)

//			fmt.Println(Spaces.Vnum)
	}
	return Spaces
}
func DescribeSpace(vnum int, Spaces []Space) {
	for i := 0; i < len(Spaces);i++ {
		if Spaces[i].Vnum == vnum {
			fmt.Println(Spaces[i].Zone)
			fmt.Println(Spaces[i].Desc)
		}
	}
}
func addPfile(play Player) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.InsertOne(context.Background(), bson.M{"name":play.Name,"title":play.Title,"inventory":play.Inventory, "equipment":play.Equipment,
						"coreboard": play.CoreBoard, "str": play.Str, "int": play.Int, "dex": play.Dex, "wis": play.Wis, "con":play.Con, "cha":play.Cha })
}
func savePfile(play Player) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("pfiles").Collection("Players")
	_, err = collection.UpdateOne(context.Background(), options.Update().SetUpsert(true), bson.M{"name":play.Name,"title":play.Title,"inventory":play.Inventory, "equipment":play.Equipment,
						"coreboard": play.CoreBoard, "str": play.Str, "int": play.Int, "dex": play.Dex, "wis": play.Wis, "con":play.Con, "cha":play.Cha })
}

//TODO make this modular
func createMobiles(name string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("mobiles").Collection("lvl1")
	_, err = collection.InsertOne(context.Background(), bson.M{"name":name,
						"str": 1, "int": 1, "dex": 1, "wis": 1, "con":1, "cha":1, "challengedice":1 })
}

//TODO make this modular
func createChat(message string, play Player) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	//process the strings
	if len(message) >= 180 {
		message = message[:180]
	}

	collection := client.Database("chat").Collection("log")
	_, err = collection.InsertOne(context.Background(), bson.M{"name":play.Name,
						"message":message, "time":time.Now(), "user":play })
	if err != nil {
		panic(err)
	}
}


func genMap(play Player, populated []Space) (Player, []Space) {
	//Create a room map
	Room := dngn.NewRoom(50, 30)
	splits := rand.Intn(75)
	Room.GenerateBSP('%', 'D', splits)
//	_, err = collection.InsertOne(context.Background(), bson.M{"room":Room})
//	if err != nil {
//		panic(err)
//	}

	populated[0].CoreBoard = ""
	play.CoreBoard = ""
	fmt.Println("Generating and populating map")
	for i := 0;i < len(Room.Data);i++ {
		if i == 0 {
			continue
		}
	//				fmt.Println(populated[0].Room.Data[populated[0].Room.Width-1][i])
			value := string(Room.Data[i])
			newValue := ""
			for s := 0;s < len(value);s++ {
				if s == 0 {
					continue
				}
				if string(value[s]) == " " {
					ChanceTreasure := "\033[48:2:200:150:0mT\033[0m"
					if rand.Intn(100) > 98 {
							newValue += ChanceTreasure
							continue
					}
					if rand.Intn(100) > 95 {
						ChanceMonster := "\033[48:2:200:50:50mM\033[0m"
						newValue += ChanceMonster
						continue
					}else {
						newValue += string(value[s])
					}
				}else {
					newValue += string(value[s])
				}

			}

			newValue = strings.ReplaceAll(newValue, "%", "\033[48:2:0:150:150m%\033[0m")
			newValue = strings.ReplaceAll(newValue, "D", "\033[38:2:200:150:150mD\033[0m")
			newValue = strings.ReplaceAll(newValue, " ", "\033[48:2:0:200:150m \033[0m")
			populated[0].CoreBoard += newValue + "\n"
			play.CoreBoard += newValue + "\n"

	}
	return play, populated
}

const (
	mapPos = "\033[0;0H"
	descPos = "\033[0;50H"
	chatStart = "\033[38:2:200:50:50m{{=\033[38:2:150:50:150m"
	chatEnd = "\033[38:2:200:50:50m=}}"
	end = "\033[0m"

)


func AssembleDescCel(room Space, row int) (string) {
	var cel string
	inWord := strings.Split(room.Desc, "\n")
	for len(strings.Split(room.Desc, "\n")) < 9 {
		room.Desc += "\n"
		inWord = strings.Split(room.Desc, "\n")
	}
	for len(inWord[0]) < 148 {
		inWord[0] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[0], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[1]) < 148 {
		inWord[1] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[1], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[2]) < 148 {
		inWord[2] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[2], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[3]) < 148 {
		inWord[3] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[3], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[4]) < 148 {
		inWord[4] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[4], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[5]) < 148 {
		inWord[5] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[5], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[6]) < 148 {
		inWord[6] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m \033[48;2;10;10;20m", inWord[6], "\033[48;2;100;5;100m \033[0m")
	for len(inWord[7]) < 148 {
		inWord[7] += " "
	}
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";51H\033[48;2;100;5;100m\033[38:2:50:0:50m@", inWord[7], "\033[48;2;100;5;100m \033[0m")

	return cel
}


func AssembleComposeCel(chatMess Chat, row int) (string, int) {
	var cel string
	inWord := chatMess.Message
	wor := ""
	word := ""
	words := ""
	if len(inWord) > 68 {
		return "DONE COMPOSTING", 0
	}
	if len(inWord) > 28 && len(inWord) > 54 {
		wor += inWord[:28]
		word += inWord[28:54]
		words += inWord[54:]
		for i := len(words); i <= 28; i++ {
			words += " "
		}
	}
	if len(inWord) > 28 && len(inWord) < 54 {
		wor += inWord[:28]
		word += inWord[28:]
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "

	}
	if len(inWord) <= 28 {
		wor = "                            "
		word += ""
		word += inWord
		for i := len(word); i <= 28; i++ {
			word += " "
		}
		words = "                            "
	}
	timeString := strings.Split(chatMess.Time.String(), " ")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", wor, "\033[48;2;10;255;20m \033[0m")
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", word, "\033[48;2;10;255;20m \033[0m"+timeString[1])
	row++
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m \033[48;2;10;10;20m", words, "\033[48;2;10;255;20m \033[0m"+timeString[0])
	row++
	namePlate := "                            "[len(chatMess.User.Name):]
	cel += fmt.Sprint("\033["+strconv.Itoa(row)+";180H\033[48;2;10;255;20m\033[38:2:50:0:50m@"+chatMess.User.Name+namePlate+"\033[48;2;10;255;20m \033[0m")

	return cel, row
	//	fmt.Println(cel)
}


func showDesc(room Space) {
	splitPos := AssembleDescCel(room, 25)
	fmt.Printf(splitPos)
	fmt.Printf("\033[38:2:140:40:140m[[")
	if room.Exits.North != 0 {
		fmt.Printf("\033[38:2:180:20:180mNorth")
	}
	if room.Exits.South != 0 {
		fmt.Printf("\033[38:2:180:20:180mSouth")
	}
	if room.Exits.East != 0 {
		fmt.Printf("\033[38:2:180:20:180mEast")
	}
	if room.Exits.West != 0 {
		fmt.Printf("\033[38:2:180:20:180mWest")
	}
	if room.Exits.NorthWest != 0 {
		fmt.Printf("\033[38:2:180:20:180mNorthWest")
	}
	if room.Exits.NorthEast != 0 {
		fmt.Printf("\033[38:2:180:20:180mNorthEast")
	}
	if room.Exits.SouthWest != 0 {
		fmt.Printf("\033[38:2:180:20:180mSouthWest")
	}
	if room.Exits.SouthEast != 0 {
		fmt.Printf("\033[38:2:180:20:180mSouthEast")
	}
	if room.Exits.Up != 0 {
		fmt.Printf("\033[38:2:180:20:180mUp")
	}
	if room.Exits.Down != 0 {
		fmt.Printf("\033[38:2:180:20:180mDown")
	}

	fmt.Printf("\033[38:2:140:40:140m]]\033[0m\033[0;0H")
	if len(room.ZonePos) >= 2 {
		drawDig(room.ZoneMap, room.ZonePos)
	}
}

func showChat(play Player) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	collection := client.Database("chat").Collection("log")
	mess, err := collection.Find(context.Background(), bson.M{})

	count := 0
	var row int
	for mess.Next(context.Background()) {
		var chatMess Chat
		err := mess.Decode(&chatMess)
		if err != nil {
			panic(err)
		}
		chatPos := fmt.Sprintf("\033["+strconv.Itoa(count+3)+";180H")
		count++
		fmt.Printf(chatPos)
		if row >= 51 {
			row = 0
		}
		message, position := AssembleComposeCel(chatMess, row)
		row = position
		fmt.Printf(message)
//		fmt.Printf(chatStart)
//		fmt.Printf(chatMess.Message + " ")
//		fmt.Printf(chatEnd)
//  	fmt.Printf(end)

	}
}
func drawDig(digFrame [][]int, zonePos []int) {
	for i := 0;i < len(digFrame);i++ {
		for c := 0;c < len(digFrame[i]);c++ {
				prn := ""
				val := fmt.Sprint(digFrame[i][c])
				if i == zonePos[0] && c == zonePos[1] {
					prn = "8"
				}
				if prn == "8" {
					fmt.Printf("\033[38:2:150:10:50m"+val+"\033[0m")
				}else if val == "1" || val == "8" {
					val = "1"
					fmt.Printf("\033[38:2:50:10:50m"+val+"\033[0m")
				}else {
						fmt.Printf(val)
				}
		}
		fmt.Println("")
	}
}

func digDug(pos []int, play Player, digFrame [][]int, digNums string, digZone string, digNum int, populated []Space) (int) {
	digVnumEnd := strings.Split(digNums, "-")[1]
	dg, digNum := initDigRoom(digFrame, digNums, digZone, play, digNum)
	play.CurrentRoom = dg
	for len(populated) <= digNum {
		populated = append(populated, dg)
	}
	populated[digNum] = dg
	dg.Vnum = digNum
	digFrame[pos[0]][pos[1]] = 8
	dg.ZonePos = dg.ZonePos[:0]
	dg.ZonePos = append(dg.ZonePos, pos[0])
	dg.ZonePos = append(dg.ZonePos, pos[1])
	fmt.Println("dug ", dg)
	drawDig(digFrame, dg.ZonePos)
	updateRoom(play, populated)
	fmt.Println("Dug ", digNum, " rooms of ", digVnumEnd)
	return digNum
}

func main() {
	//TODO Get the Spaces that are already loaded in the database and skip
	//if vnum is taken
	//Get the flags passed in
	var populated []Space
	var play Player
	//Make this relate to character level
	var dug []Space
	if len(os.Args) > 1 {
		if os.Args[1] == "--init" {
			//TODO testing suite - one test will be randomly generating 10,000 Spaces
			//and seeing if the system can take it
			descString := "The absence of light is blinding.\nThree large telephone poles illuminate a small square."
			for len(strings.Split(descString, "\n")) < 8 {
				descString += "\n"
			}
			InitZoneSpaces("0-5", "The Void", descString)
			descString = "I wonder what day is recycling day.\nEven the gods create trash."
			for len(strings.Split(descString, "\n")) < 8 {
				descString += "\n"
			}
			InitZoneSpaces("5-15", "Midgaard", descString)
			populated = PopulateAreas()
			play = InitPlayer("FSM")
			addPfile(play)
			createMobiles("Noodles")
		}else if os.Args[1] == "--user" {
			//Continue on
			populated = PopulateAreas()
			play = InitPlayer("Wallace")
			savePfile(play)
			fmt.Println("In client loop")
			fmt.Printf("\033[51;0H")
		}else if os.Args[1] == "--builder" {
			//Continue on
			populated = PopulateAreas()
			play = InitPlayer("FlyingSpaghettiMonster")
			savePfile(play)

			fmt.Println("Builder log-in")

			fmt.Printf("\033[51;0H")
		}else {
			fmt.Println("Unrecognized flag")
			os.Exit(1)
		}
	} else {
		fmt.Println("Use --init to build and launch the world, --user to just connect.")
		fmt.Println("--builder for a building session")
		os.Exit(1)
	}


	//Game loop
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(){
		clear()
		savePfile(play)
		input := scanner.Text()
		//Save pfile first
		save := false
		if strings.HasPrefix(input, "dig") {
			var digFrame [][]int
			for i := 0;i < 50;i++ {
				Frame := make([]int, 50)
				digFrame = append(digFrame, Frame)
			}
			fmt.Println("\033[38:2:255:0:0m", len(digFrame), "\033[0m")

			//Make a bar that fills with how many rooms you dig
			pos := make([]int, 2)
			pos[0] = 25
			pos[1] = 25
			if len(strings.Split(input, " ")) == 4 {
				digZone := strings.Split(input, " ")[1]
				digVnumStart := strings.Split(input, " ")[2]
				digVnumEnd := strings.Split(input, " ")[3]

				//Error was nil so start the digging protocol
				save = false
				dug = dug[:0]

				digNums := digVnumStart + "-" + digVnumEnd
				digNum, err := strconv.Atoi(digVnumStart)
				if err != nil {
					panic(err)
				}
				DIG:
				for scanner.Scan() {
					input = scanner.Text()
					inp, err := strconv.Atoi(input)
					if err != nil {
						fmt.Sprint("\033[0;0HAlphabetic code entry found")
						switch input {
						case "update zonemap":
							updateZoneMap(play, populated)
						case "edit desc":
							//desc
							//room has to exist before we edit it
							digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							//dignum shouldn't change because we're editing the same room

							play.CurrentRoom.Desc = ""
							fmt.Println("Enter the room's new description, enter for a new line, @ on a new line to end.")
							descScanner := bufio.NewScanner(os.Stdin)
							DESC:
							for descScanner.Scan() {
								if descScanner.Text() == "@" || len(strings.Split(populated[play.CurrentRoom.Vnum].Desc, "\n")) < 8 {
									if descScanner.Text() == "@" {
										for len(strings.Split(populated[play.CurrentRoom.Vnum].Desc, "\n")) < 8 {
											populated[play.CurrentRoom.Vnum].Desc += "\n"
										}
									}
									populated[play.CurrentRoom.Vnum].Desc = play.CurrentRoom.Desc
									break DESC
								}else {
									play.CurrentRoom.Desc += descScanner.Text() + "\n"
								}
							}
							client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
							if err != nil {
								panic(err)
							}
							ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
							err = client.Connect(ctx)
							if err != nil {
								panic(err)
							}
							filter := bson.M{"vnum": play.CurrentRoom.Vnum}
							collection := client.Database("zones").Collection("Spaces")
							update := bson.M{"$set": bson.M{"vnums":populated[play.CurrentRoom.Vnum].Vnums,
								 "desc":populated[play.CurrentRoom.Vnum].Desc,"exits": populated[play.CurrentRoom.Vnum].Exits,
									 "altered": true }}

							result, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
							if err != nil {
								panic(err)
							}
							fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
						case "edit title":
							//room title
						case "edit mobiles":
							//mobiles
						case "edit items":
							//items
						default:
							fmt.Println("I don't understand")
						}

						err = nil
					}
					//Set up the whole keypad for "digging"
					switch inp {
					case 1101:
						save = false
						break DIG
					case 1111:
						save = true
						break DIG
					case 1:
						//Sw

						if digFrame[pos[0]+1][pos[1]-1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] += 1
							pos[1] -= 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 2:
						//S
						if digFrame[pos[0]+1][pos[1]] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] += 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 3:
						//Se
						if digFrame[pos[0]+1][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] += 1
							pos[1] += 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 4:
						//W
						if digFrame[pos[0]][pos[1]-1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[1] -= 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
							}
					case 5:
						//TODO, make a selector for which level is shown
						//Down
						dg, _ := initDigRoom(digFrame, digNums, digZone, play, digNum)

						digNum = dg.Vnum
						fmt.Println("dug ", dg)

						play.CurrentRoom = dg
						populated[dg.Vnum] = dg

						save = true
					case 6:
						//E
						if digFrame[pos[0]][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[1] += 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 7:
						//Nw
						if digFrame[pos[0]-1][pos[1]-1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							pos[1] -= 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 8:
						//N
						if digFrame[pos[0]-1][pos[1]] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 9:
						//Ne
						if digFrame[pos[0]-1][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							pos[1] += 1
							digNum = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					default:
						if len(play.CurrentRoom.ZonePos) >= 2 {
							drawDig(digFrame, play.CurrentRoom.ZonePos)
						}
						fmt.Println("Dug ", digNum, " rooms of ", digVnumEnd)
					}
				}


			}
			if save {
				client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
				if err != nil {
					panic(err)
				}
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				err = client.Connect(ctx)
				if err != nil {
					panic(err)
				}

				file, err := os.Create("dat/zone.bson")
				if err != nil {
					panic(err)
				}
				defer file.Close()
				writer := bufio.NewWriter(file)
				fmt.Println("\033[38:2:200:50:50mUpdating the zone with final map.\033[0m")
				updateZoneMap(play, populated)
				fmt.Println("Dumping the area list to dat/zone.bson")
				for i := 0;i < len(populated);i++ {
					marshalledBson, err := bson.Marshal(populated[i])
					if err != nil {
						panic(err)
					}
					writer.Write(marshalledBson)
					writer.Flush()
				}
			}

			}



		if input == "edit desc"{

			play.CurrentRoom.Desc = ""
			fmt.Println("Enter the room's new description, enter for a new line, @ on a new line to end.")
			descScanner := bufio.NewScanner(os.Stdin)
			DESCREG:
			for descScanner.Scan() {
				if descScanner.Text() == "@" || len(strings.Split(populated[play.CurrentRoom.Vnum].Desc, "\n")) < 8 {
					if descScanner.Text() == "@" {
						for len(strings.Split(populated[play.CurrentRoom.Vnum].Desc, "\n")) < 8 {
							populated[play.CurrentRoom.Vnum].Desc += "\n"
						}
					}
					populated[play.CurrentRoom.Vnum].Desc = play.CurrentRoom.Desc
					break DESCREG
				}else {
					play.CurrentRoom.Desc += descScanner.Text() + "\n"
				}
			}

			client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
			if err != nil {
				panic(err)
			}
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			err = client.Connect(ctx)
			if err != nil {
				panic(err)
			}
			filter := bson.M{"vnum": play.CurrentRoom.Vnum}
			collection := client.Database("zones").Collection("Spaces")
			update := bson.M{"$set": bson.M{"vnums":populated[play.CurrentRoom.Vnum].Vnums,
				 "desc":populated[play.CurrentRoom.Vnum].Desc,"exits": populated[play.CurrentRoom.Vnum].Exits,
					 "altered": true }}

			result, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
			if err != nil {
				panic(err)
			}
			fmt.Println("\033[38:2:255:0:0m", result, "\033[0m")
		}



		if input == "quit" {
			fmt.Println("Bai!")
			os.Exit(1)
		}
		if strings.HasPrefix(input, "ooc") {
			createChat(input[3:], play)
			showChat(play)
		}

		if input == "look" {
			fmt.Sprintf("Current room is ", play.CurrentRoom)
			showDesc(play.CurrentRoom)
		}
		if strings.Contains(input, "gen coreboard") {
			play, populated = genMap(play, populated)
		}
		if strings.Contains(input, "open map") {
			//// TODO:
			//This
		}

		if strings.Contains(input, "open coreboard") {
			fmt.Printf(mapPos)
			fmt.Print(populated[0].CoreBoard)
		}
		if strings.HasPrefix(input, "view from") {
			splitCommand := strings.Split(input, "from")
			stripped := strings.TrimSpace(splitCommand[1])
			vnumLook, err := strconv.Atoi(stripped)
			if err != nil {
				fmt.Println("Error converting a stripped string")
			}
			DescribeSpace(vnumLook, populated)
		}
		if strings.HasPrefix(input, "go to") {
			splitCommand := strings.Split(input, "to")
			stripped := strings.TrimSpace(splitCommand[1])
			inp, err := strconv.Atoi(stripped)
			if err != nil {
				fmt.Println("Error converting a stripped string")
			}
			for i := 0;i < len(populated);i++ {
				if inp == populated[i].Vnum {
					play.CurrentRoom = populated[i]
					fmt.Print(populated[i].Vnum, populated[i].Vnums, populated[i].Zone)
					showDesc(play.CurrentRoom)
					fmt.Printf("\033[0;0H\033[38:2:0:255:0mPASS\033[0m")
					break
				}else {
					fmt.Printf("\033[0;0H\033[38:2:255:0:0mERROR\033[0m")
				}
			}
		}
		if input == "score" {
			DescribePlayer(play)
		}
		//Reset the input to a standardized place
		showDesc(play.CurrentRoom)
		showChat(play)
		fmt.Printf("\033[51;0H")
	}
//	res, err := collection.InsertOne(context.Background(), bson.M{"Noun":"x"})
//	res, err = collection.InsertOne(context.Background(), bson.M{"Verb":"+"})
//	res, err = collection.InsertOne(context.Background(), bson.M{"ProperNoun":"y"})

}
