package main

import (
  "fmt"
  "os"
	"time"
	"context"
  "strconv"
  "bufio"
  "math/rand"
  "strings"
  zmq "github.com/pebbe/zmq4"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)
func ShowOoc(response *zmq.Socket, play Player) string {
  input := "+++"
  input = play.Name+input
  //createChat(input[3:], play)
  //todo
  response.Recv(0)
  _, err := response.Send(input, 0)
  if err != nil {
    panic(err)
  }
  chat, err := response.Recv(0)
  if err != nil {
    panic(err)
  }
  out := fmt.Sprint(chat)
	return out
}
func JackIn(in chan bool) error {
  fmt.Printf("\033[10;28H\033[0m")
  fmt.Printf("\033[11;28H \033[48;2;10;255;20m\033[38;2;10;10;255m         LOGIN         \033[0m")
  fmt.Printf("\033[12;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[13;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mUSER                \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[14;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[15;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[16;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mPASSWORD            \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[17;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[18;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[19;28H \033[48;2;10;255;20m                       \033[0m")
  fmt.Printf("\033[17;32H")
out := ""
row := 0
for i := 0;i < 52;i++ {
  for count := 0;count < 250;count++ {
    select {
    case notConn := <- in:
      clearDirty()
      if notConn == false {

        return nil
      }

    default:
          if rand.Intn(45) > 35 {
            randPosX := strconv.Itoa(rand.Intn(200))
            randPosY := strconv.Itoa(rand.Intn(52))
            out += "\033["+randPosY+";"+randPosX+"H\033[48:2:250:250:250m \033[0m"
          }else {
            randPosX := strconv.Itoa(rand.Intn(200))
            randPosY := strconv.Itoa(rand.Intn(52))
            out += "\033["+randPosY+";"+randPosX+"H\033[48:2:25:35:25m \033[0m"
          }
          row++

          time.Sleep(10*time.Millisecond)
          fmt.Print(out)

    }
  }

  }
  return nil
}


func LoginSC() (string, string){
  clearDirty()
  loginScanner := bufio.NewScanner(os.Stdin)
  user := ""
  fmt.Printf("\033[10;28H\033[0m")
  fmt.Printf("\033[11;28H \033[48;2;10;255;20m\033[38;2;10;10;255m         LOGIN         \033[0m")
  fmt.Printf("\033[12;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[13;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mUSER                \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[14;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[15;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[16;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mPASSWORD            \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[17;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[18;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
  fmt.Printf("\033[19;28H \033[48;2;10;255;20m                       \033[0m")
  fmt.Printf("\033[14;32H" + user + "\033[0m")

  loginScanner.Scan()

  user = loginScanner.Text()
  fmt.Printf("\033[17;32H")
  for {

		fmt.Printf("\033[10;28H\033[0m")
		fmt.Printf("\033[11;28H \033[48;2;10;255;20m\033[38;2;10;10;255m         LOGIN         \033[0m")
		fmt.Printf("\033[12;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[13;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mUSER                \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[14;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[15;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[16;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   \033[38;2;10;200;150mPASSWORD            \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[17;28H\033[48;2;10;255;20m \033[48;2;10;10;20m   ________________    \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[18;28H\033[48;2;10;255;20m \033[48;2;10;10;20m                       \033[48;2;10;255;20m \033[0m")
		fmt.Printf("\033[19;28H \033[48;2;10;255;20m                       \033[0m")
		fmt.Printf("\033[14;32H" + user + "\033[0m")
    fmt.Printf("\033[17;32H")
		loginScanner.Scan()
    pword := loginScanner.Text()
    //clearDirty()
    //Only use clearDirty at major intersections, it will cause flicker
		return user, pword

  }
  return "", ""
}
func clearDirty() {
  for i := 0;i < 255;i++ {
    fmt.Println("")
  }
}

func showBattle(damMsg []string) string {
	out := ""
  clear := false
  outClear := ""
  for i := 0;i < len(damMsg);i++ {
    if len(damMsg) > 17 {
        damMsg = damMsg[17:]
        //clearDirty()
        //reset()
        i = 0
        clear = true
    }
    if clear {
      for c := 0;c < 17;c++ {
        outClear += fmt.Sprint("\033["+strconv.Itoa(i)+";53H                                                                            ")
        if c == 16 {
            fmt.Print(outClear)
            clear = false
        }
      }

    }else {
      out += fmt.Sprint("\033["+strconv.Itoa(i+1)+";53H"+damMsg[i]+"\033["+strconv.Itoa(i+2)+";53H                                                                    ")

    }
  }
  fmt.Print(outClear)
	return out
}

func clearCmd() {
		fmt.Print(cmdPos+"                                                                                                                                                                                   ")
		fmt.Print("\033[52;0H                                                                                                                                                                                   ")
		fmt.Print("\033[53;0H                                                                                                                                                                                   ")
		fmt.Print("\033[54;0H                                                                                                                                                                                   ")
		fmt.Print("\033[55;0H                                                                                                                                                                                   ")
		fmt.Print("\033[56;0H                                                                                                                                                                                   ")
		fmt.Print(cmdPos)
}

func showCoreMobs(play Player) (Player, string) {
  core := ""
	out := ""
  coreSplit := strings.Split(play.PlainCoreBoard, "\n")
  for i := 0;i < len(coreSplit);i++ {
    for r := 0;r < len(coreSplit[i]);r++ {
      if coreSplit[i][r] == 'M' {
          for bat := 0;bat < len(play.Fights.Oppose);bat++ {
            if play.Fights.Oppose[bat].MaxRezz <= 0 && play.Fights.Oppose[bat].X == r && play.Fights.Oppose[bat].Y == i{
              //fmt.Println("ONE DOWN AT"+strconv.Itoa(play.Fights.Oppose[bat].X)+":"+strconv.Itoa(play.Fights.Oppose[bat].Y))
              play.Fights.Oppose[bat].Char = fmt.Sprint("\033[48;2;5;0;150m\033["+strconv.Itoa(play.Fights.Oppose[bat].Y+20)+";"+strconv.Itoa(play.Fights.Oppose[bat].X+54)+"H\033[48:2:175:0:0mC\033[0m")
    //          core += play.Fights.Oppose[bat].Char
  //            play.TargetLong = "C"
              break
            }else {
//              play.TargetLong = string(coreSplit[i][r])
              core += fmt.Sprint("\033["+strconv.Itoa(i+20)+";"+strconv.Itoa(r+54)+"H\033[48:2:175:0:150m"+string(play.Fights.Oppose[bat].Char)+"\033[0m")

            }

        }


        }
      }
    }

  out += fmt.Sprint(core)
  return play, out
}

func showCoreBoard(play Player) string {
  core := ""
	out := ""
  coreSplit := strings.Split(play.CoreBoard, "\n")
  for i := 0;i < len(coreSplit);i++ {
      core += fmt.Sprint("\033[",strconv.Itoa(i+20),";54H",coreSplit[i])
    }
    out = fmt.Sprint(core)
		return out
}
func clearCoreBoard(play Player) {
  coreSplit := strings.Split(play.CoreBoard, "\n")
  //This needs to be made dynamic for when we adjust the view. for now it's fine
  coreSpace := "                          "
  for i := 0;i < len(coreSplit);i++ {
    core := fmt.Sprint("\033[",strconv.Itoa(i+20),";52H ")

    fmt.Print(core+coreSpace)
  }
}
func DescribeSpace(vnum int, Spaces []Space) string {
	out := ""
	for i := 0; i < len(Spaces);i++ {
		if Spaces[i].Vnum == vnum {
			out += fmt.Sprint(Spaces[i].Zone)
			out += fmt.Sprint(Spaces[i].Desc)
		}
	}
	return out
}

func showDesc(room Space) string {
	out := ""
	splitPos := AssembleDescCel(room, 25)
	out += fmt.Sprint(splitPos)
	out += fmt.Sprint("\033[38:2:140:40:140m[[")
	if room.Exits.North != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mNorth")
	}
	if room.Exits.South != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mSouth")
	}
	if room.Exits.East != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mEast")
	}
	if room.Exits.West != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mWest")
	}
	if room.Exits.NorthWest != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mNorthWest")
	}
	if room.Exits.NorthEast != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mNorthEast")
	}
	if room.Exits.SouthWest != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mSouthWest")
	}
	if room.Exits.SouthEast != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mSouthEast")
	}
	if room.Exits.Up != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mUp")
	}
	if room.Exits.Down != 0 {
		out += fmt.Sprint("\033[38:2:180:20:180mDown")
	}

	out += fmt.Sprint("\033[38:2:140:40:140m]]\033[0m\033[0;0H")
	if len(room.ZonePos) >= 2 {
		out += drawDig(room.ZoneMap, room.ZonePos)
	}
	return out
}
func showWho(play Player) []string {
  var whoList []string
  client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
  filter := bson.M{}
  collection := client.Database("who").Collection("players")
  findOptions := options.Find()
  findOptions.SetLimit(1000)
  results, err := collection.Find(context.Background(), filter, findOptions)
  if err != nil {
    panic(err)
  }
  for results.Next(context.Background()) {
    var signedIn SignIn
    err := results.Decode(&signedIn)
    if err != nil {
      panic(err)
    }
    whoList = append(whoList, signedIn.Payload.Name)
  }
  return whoList
}
func updateWho(play Player, in bool) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
  filter := bson.M{"event":"players/sign-in","payload":bson.M{"name":play.Name}}
  update := bson.M{"event":"players/sign-in","ref":UIDMaker(), "payload":bson.M{"name":play.Name}}
	collection := client.Database("who").Collection("players")
	findOptions := options.Find()
  findOptions.SetLimit(1000)
  if !in {
    _, err := collection.DeleteOne(context.Background(), filter)
  	if err != nil {
  		panic(err)
  	}
  }else {
    _, err := collection.InsertOne(context.Background(), update)
    if err != nil {
      panic(err)
    }
  }

}
func showChat(play Player) (int, string) {
	out := ""
  countChat := 0
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
    count++
		var chatMess Chat
		err := mess.Decode(&chatMess)
		if err != nil {
			panic(err)
		}
		chatPos := fmt.Sprintf("\033["+strconv.Itoa(count+3)+";180H")
		countChat++
		fmt.Printf(chatPos)
		if row >= 51 {
			row = 0
		}
		message, position := AssembleComposeCel(chatMess, row)
		row = position
		out += fmt.Sprint(message)
//		fmt.Printf(chatStart)
//		fmt.Printf(chatMess.Message + " ")
//		fmt.Printf(chatEnd)
//  	fmt.Printf(end)

	}
	return countChat, out
}
func drawDig(digFrame [][]int, zonePos []int) string {
	out := ""
	for i := 0;i < len(digFrame);i++ {
		out += fmt.Sprint("\033[48;2;10;255;20m \033[0m")
		for c := 0;c < len(digFrame[i]);c++ {
				prn := ""
				val := fmt.Sprint(digFrame[i][c])
				if i == zonePos[0] && c == zonePos[1] {
					prn = "8"
				}
				if prn == "8" {
					out += fmt.Sprint("\033[38:2:150:10:50m"+val+"\033[0m")
				}else if val == "1" || val == "8" {
					val = "1"
					out += fmt.Sprint("\033[38:2:50:10:50m"+val+"\033[0m")
				}else if c == 0 || c == len(digFrame[i])-1 || i == 0 || i == len(digFrame)-1{
          out += fmt.Sprint("\033[48;2;10;255;20m \033[0m")
        }else {
						out += fmt.Sprint(val)
				}
		}
		out += fmt.Sprintln("\033[48;2;10;255;20m \033[0m")
	}
	return out
}

func DescribePlayer(play Player) string {
	out := ""
  ratio := ""
  count := 18
  for   rezz := 0;rezz < play.Rezz;rezz++ {

    ratio += "\033["+strconv.Itoa(count+30)+";25H\033[48:2:175:50:50m \033[0m\n"
    count--
  }
  for count > 0 {
      ratio += "\033["+strconv.Itoa(count+30)+";25H\033[48:2:15:50:50m \033[0m\n"

    count--
  }

  ratio += "\033[31;24H+++\n"
  ratio += "\033[49;24H+++"
  hp := ratio
  count = 18
  ratio = ""
  for tech := 0;tech < play.Tech;tech++ {
    ratio += "\033["+strconv.Itoa(count+30)+";31H\033[48:2:75:150:50m \033[0m\n"
    count--
  }
  for count > 0 {
      ratio += "\033["+strconv.Itoa(count+30)+";31H\033[48:2:15:50:50m \033[0m\n"
      count--
  }
  ratio += "\033[31;30H===\n"
  ratio += "\033[49;30H==="
  techShow := ratio
  out += fmt.Sprint(techShow)
  out += fmt.Sprint(hp)
	out += fmt.Sprint("\033[40;0H")
	out += fmt.Sprintln("======================")
	out += fmt.Sprintln("\033[38:2:0:200:0mStrength     :\033[0m", play.Str)
	out += fmt.Sprintln("\033[38:2:0:200:0mIntelligence :\033[0m", play.Int)
	out += fmt.Sprintln("\033[38:2:0:200:0mDexterity    :\033[0m", play.Dex)
	out += fmt.Sprintln("\033[38:2:0:200:0mWisdom       :\033[0m", play.Wis)
	out += fmt.Sprintln("\033[38:2:0:200:0mConstitution :\033[0m", play.Con)
	out += fmt.Sprintln("\033[38:2:0:200:0mCharisma     :\033[0m", play.Cha)
	out += fmt.Sprintln("======================")
	return out
}
