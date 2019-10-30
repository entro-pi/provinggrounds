package main

import (
	"bufio"
	"os"
	"context"
	"time"
	"fmt"
	"strconv"
	"strings"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/SolarLune/dngn"
)



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

type Player struct {
	Name string
	Title string
	Inventory []int
	Equipment []int
	CoreBoard string
	PlainCoreBoard string
	CurrentRoom Space

	Rezz int
	Tech int

	Str int
	Int int
	Dex int
	Wis int
	Con int
	Cha int
}




const (
	cmdPos = "\033[51;0H"
	mapPos = "\033[1;51H"
	descPos = "\033[0;50H"
	chatStart = "\033[38:2:200:50:50m{{=\033[38:2:150:50:150m"
	chatEnd = "\033[38:2:200:50:50m=}}"
	end = "\033[0m"

)



func main() {
	//TODO Get the Spaces that are already loaded in the database and skip
	//if vnum is taken
	//Get the flags passed in
	var populated []Space
	var play Player
	//Make this relate to character level
	var dug []Space
	coreShow := false
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
			fmt.Println("\033[38:2:0:250:0mAll tests passed and world has been initialzed\n\033[0mYou may now start with --login.")
			os.Exit(1)
		}else if os.Args[1] == "--guest" {
			//Continue on
			populated = PopulateAreas()
			play = InitPlayer("Wallace")
			savePfile(play)
			fmt.Println("In client loop")
			fmt.Printf("\033[51;0H")
		}else if os.Args[1] == "--login" {
			//Continue on
			user, pword := LoginSC()

			populated = PopulateAreas()
			play = InitPlayer(user)
			//just hang on to the password for now
			fmt.Sprint(pword)
			savePfile(play)
			fmt.Println("In client loop")
			input := "go to 1"
			//this is pretty incomprehensible
			//TODO
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
					DescribePlayer(play)
					fmt.Printf("\033[0;0H\033[38:2:0:255:0mPASS\033[0m")
					break
				}else {
					fmt.Printf("\033[0;0H\033[38:2:255:0:0mERROR\033[0m")
				}
			}
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
	firstDig := false
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan(){
		clearCmd()
		savePfile(play)
		input := scanner.Text()
		//Save pfile first
		save := false
		if strings.HasPrefix(input, "dig") {
			if strings.Split(input, " ")[1] == "new" {
				firstDig = true
			}else {
				firstDig = false
			}
			if firstDig {
				fmt.Println("Now specify the zone name and vnums required")
				fmt.Println("as in, \"dig zem 0 15\"")
				scanner.Scan()
				input = scanner.Text()
			}
			var digFrame [][]int
			for i := 0;i < 30;i++ {
				Frame := make([]int, 50)
				digFrame = append(digFrame, Frame)
			}

			fmt.Println("\033[38:2:255:0:0m", len(digFrame), "\033[0m")

			//Make a bar that fills with how many rooms you dig

			pos := make([]int, 2)

			if firstDig {
				pos[0] = 25
				pos[1] = 25
			}else {
				pos[0] = play.CurrentRoom.ZonePos[0]
				pos[1] = play.CurrentRoom.ZonePos[1]

			}

			if len(strings.Split(input, " ")) == 4 {
				digZone := strings.Split(input, " ")[1]
				digVnumStart := strings.Split(input, " ")[2]
				digVnumEnd := strings.Split(input, " ")[3]

				//Error was nil so start the digging protocol
				save = false
				dug = dug[:0]

				digNums := digVnumStart + "-" + digVnumEnd
				toDig := PopulateAreaBuild(digNums)
				for i := 0;i < len(toDig);i++ {
					populated = append(populated, toDig[i])

				}

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
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 2:
						//S
						if digFrame[pos[0]+1][pos[1]] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] += 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 3:
						//Se
						if digFrame[pos[0]+1][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] += 1
							pos[1] += 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 4:
						//W
						if digFrame[pos[0]][pos[1]-1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[1] -= 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
							}
					case 5:
						//TODO, make a selector for which level is shown
						//Down

						save = true
					case 6:
						//E
						if digFrame[pos[0]][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[1] += 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 7:
						//Nw
						if digFrame[pos[0]-1][pos[1]-1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							pos[1] -= 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 8:
						//N
						if digFrame[pos[0]-1][pos[1]] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
							play.CurrentRoom.Vnum = digNum
						}
					case 9:
						//Ne
						if digFrame[pos[0]-1][pos[1]+1] != 1 {
							digFrame[pos[0]][pos[1]] = 1
							pos[0] -= 1
							pos[1] += 1
							digNum, play.CurrentRoom = digDug(pos, play, digFrame, digNums, digZone, digNum, populated)
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


		//secondary commands
		if input == "targeting computer" {
			fmt.Print("Input co-ordinates in the form of AA AB AC etc..")
			err := target(play, populated)
			if err != nil {
				panic(err)
			}
		}
		if input == "show room vnum" {
			fmt.Print("\033[38;2;150;0;150mROOM VNUM :"+strconv.Itoa(play.CurrentRoom.Vnum)+"\033[0m")
		}
		if input == "dam rezz" {
			play.Rezz -= 5
		}
		if input == "dam tech" {
			play.Tech -= 6
		}
		if input == "heal" {
			play.Rezz = 17
			play.Tech = 17
		}
		if input == "show zone info" {
			fmt.Println("\033[38;2;150;0;150mZONE NAME :"+play.CurrentRoom.Zone+"\033[0m")
			fmt.Print("\033[38;2;150;0;150mZONE VNUMS :"+play.CurrentRoom.Vnums+"\033[0m")
		}
		if input == "edit desc"{

			play.CurrentRoom.Desc = ""
			fmt.Println("Enter the room's new description, enter for a new line, @ on a new line to end.")
			descScanner := bufio.NewScanner(os.Stdin)
			DESCREG:
			for descScanner.Scan() {
				if descScanner.Text() == "@" {
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
			populated = PopulateAreas()
		}



		if input == "quit" {
			fmt.Println("Bai!")
			os.Exit(1)
		}
		if strings.HasPrefix(input, "ooc") {
			createChat(input[3:], play)
			showChat(play)
		}
		if input == "blit" {
			clearDirty()
		}
		if input == "count keys" {
			countKeys()
			showDesc(play.CurrentRoom)
		}
		if strings.HasPrefix(input, "merge") {
			fmt.Println("Merging area zone map data")
			split := strings.Split(input, " ")
			sourceName, destName := split[1], split[2]
			var sourceDat [][]int
			var destDat [][]int
			for i := 0;i < len(populated);i++ {
				if populated[i].Zone == sourceName {
					sourceDat = populated[i].ZoneMap
				}
			}
			for i := 0;i < len(populated);i++ {
				if populated[i].Zone == destName {
					destDat = populated[i].ZoneMap
				}
			}
			zoneDat := mergeMaps(sourceDat, destDat)
			populated[play.CurrentRoom.Vnum].ZoneMap = zoneDat
			play.CurrentRoom.ZoneMap = zoneDat
			play.CurrentRoom.Zone = sourceName
			updateZoneMap(play, populated)
			play.CurrentRoom.Zone = destName
			updateZoneMap(play, populated)
		}
		if input == "update zonemap" {
			updateZoneMap(play, populated)
		}
		if input == "look" {
			fmt.Sprintf("Current room is ", play.CurrentRoom)
			showDesc(play.CurrentRoom)
			DescribePlayer(play)
		}
		if strings.Contains(input, "gen coreboard") {
			//TODO make this so one doesn't loose the
			//old coreboard, or convert it to xp, i dunno
			play.CoreBoard, play = genCoreBoard(play, populated)
		}
		if strings.Contains(input, "open map") {
			//// TODO:
			//This
		}
		if strings.Contains(input, "close coreboard") {
			coreShow = false
		}
		if strings.Contains(input, "lock coreboard") {
//			fmt.Printf(mapPos)
				showCoreBoard(play)
				coreShow = true
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
					DescribePlayer(play)
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
		DescribePlayer(play)
		showChat(play)
		if coreShow {
			showCoreBoard(play)
		}else {
			clearCoreBoard(play)
		}
		fmt.Printf("\033[51;0H")
	}
//	res, err := collection.InsertOne(context.Background(), bson.M{"Noun":"x"})
//	res, err = collection.InsertOne(context.Background(), bson.M{"Verb":"+"})
//	res, err = collection.InsertOne(context.Background(), bson.M{"ProperNoun":"y"})

}
