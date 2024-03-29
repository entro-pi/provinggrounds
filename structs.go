package main

import (
	"github.com/SolarLune/dngn"
	"time"
)

type Class struct {
	Level float64
	Name string
	Skills []Skill
	Spells []Spell
}

type Spell struct {
	TechUsage int
	Usage rune
	Level int
	Consumed bool
	Name string
}

type Skill struct {
	Name string
	DamType string
	Level int
	Usage rune
	Dam int
}

type StatusPayload struct {
	Game string
	Players []string
}

type Status struct {
	Event string
	Ref string
	Payload StatusPayload
}


type SignOutPayload struct {
	Name string
	Game string
}

type SignOut struct {
	Event string
	Ref string
	Payload SignOutPayload
}

type SignIn struct {
	Event string
	Ref string
	Payload SignInPayload
}
type InventoryItem struct {
	Item Object
	Number int
	Amount int
}
type SignInPayload struct {
	Name string
	Game string
}
type Object struct {
	Name string
	LongName string
	Vnum int
	Zone string
	Owner Player
	Value int
	Slot int
	X int
	Y int
	Owned bool
	Amount int
	Number int
}

type BroadcastPayload struct {
  Channel string
  Message string
  Game string
  Name string
	Row int
	Col int
	Selected bool
	BigMessage string
	ID int
	Transaction OnlineTransaction
	Store OnlineStore
}
type Butler struct {
	Employer string
	Funds Account
	ToBuy Object
}
type OnlineStore struct {
	Owner string
	Float float64
	Inventory []OnlineTransaction
}

type OnlineTransaction struct {
	ItemHash string
	Item Object
	Sold bool
	To Account
	Price float64
}
type Bank struct {
	Clientele Client
	Owner string
}
type Client struct {
		User string
		TotalAmount float64
		Accounts []float64
}
type Account struct {
	Owner string
	Income OnlineStore
	Amount float64
}

type Broadcast struct {
    Event string
    Ref string
    Payload BroadcastPayload
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
	MobilesInRoom []Mobile
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
type EquipmentItem struct {
	Item Object
}

type Player struct {
	Name string
	Title string
	InventoryStore InventoryBank
	Inventory []Object
	Equipment []int
	Equipped []EquipmentItem
	CoreBoard string
	PlainCoreBoard string
	CurrentRoom Space
	PlayerHash string
	Classes []Class
	Target string
	TargetLong string
	//ToBuy will have to either be an ItemHash
	//Or a vnum
	ToBuy int
	BankAccount Account
	TarX int
	TarY int
	OldX int
	OldY int
	CPU string
	CoreShow bool
	Channels []string
	Battling bool
	Profile string
	Session string

	Slain int
	Hoarded int


	MaxRezz int
	Rezz int
	Tech int
	Fights Fight
	Won int
	Found int

	Str int
	Int int
	Dex int
	Wis int
	Con int
	Cha int
}
type Fight struct {
	Oppose []Mobile
	Former []Player
	Treasure []Object
}
type Mobile struct {
	Name string
	LongName string
	ItemSpawn []int
	Rep string
	MaxRezz int
	Rezz int
	Tech int
	Aggro bool
	Align int
	Vnum int
	X int
	Y int
	Char string
}
