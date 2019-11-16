package main

type InventoryBank struct {

  SlotOne InventoryItem
  SlotOneAmount int

  SlotTwo InventoryItem
  SlotTwoAmount int

  SlotThree InventoryItem
  SlotThreeAmount int

  SlotFour InventoryItem
  SlotFourAmount int

  SlotFive InventoryItem
  SlotFiveAmount int

  SlotSix InventoryItem
  SlotSixAmount int

  SlotSeven InventoryItem
  SlotSevenAmount int

  SlotEight InventoryItem
  SlotEightAmount int

  SlotNine InventoryItem
  SlotNineAmount int

  SlotTen InventoryItem
  SlotTenAmount int
}

//Todo, make these receivers on a type
func initInv(play Player) Player {
  for i := 0;i < play.InventoryStore.SlotOneAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotOne)
    }
  }

  for i := 0;i < play.InventoryStore.SlotTwoAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotTwo)
    }
  }

  for i := 0;i < play.InventoryStore.SlotThreeAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotThree)
    }
  }

  for i := 0;i < play.InventoryStore.SlotFourAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotFour)
    }
  }

  for i := 0;i < play.InventoryStore.SlotFiveAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotFive)
    }
  }

  for i := 0;i < play.InventoryStore.SlotSixAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotSix)
    }
  }

  for i := 0;i < play.InventoryStore.SlotSevenAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotSeven)
    }
  }

  for i := 0;i < play.InventoryStore.SlotEightAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotEight)
    }
  }

  for i := 0;i < play.InventoryStore.SlotNineAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotNine)
    }
  }

  for i := 0;i < play.InventoryStore.SlotTenAmount;i++ {
    if i >= 1 {
      play.Inventory[len(play.Inventory)-1].Amount++
    }else {
      play.Inventory = append(play.Inventory, play.InventoryStore.SlotTen)
    }
  }

  return play
}
