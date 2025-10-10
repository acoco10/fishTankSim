//go:build old

package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"log"
)

func (wb *WhiteBoardSprite) UpdateOld() {

	sp := wb.UiSprite.Sprite
	wb.graphicManager.UpdateTimers()

	if wb.wbState == WritingTask {
		if wb.CheckDoneWriting() {
			wb.wbState = IdleWB
			if wb.taskState == ClickToCompleteAfterWriting {
				AddClickme(wb)
			}
		}
	}

	if sp.SpriteHovered() && ClickCheck() && wb.taskJustCompleted {
		if !wb.windowOpen {
			wb.tasksCompletedToday++
			wb.taskJustCompleted = false
			ev2 := tasks.TaskCompleted{Task: *wb.currentTask}
			wb.UiSprite.EventHub.Publish(ev2)
			wb.UiSprite.Sprite.DOptsUpdaterTag = ""
			wb.UiSprite.Sprite.Shader = nil
		}
	}

	if wb.allTasksCompleted {
		wb.allTasksCompleted = false
		wb.timers[AllTasksTimer].TurnOn()
		AddClickme(wb)
		//add click me for initiating erase
	}

	if wb.allTasksBufferDone && sp.SpriteHovered() && ClickCheck() {
		wb.allTasksBufferDone = false
		wb.UiSprite.Sprite.DOptsUpdaterTag = ""
		wb.UiSprite.Sprite.Shader = nil
		wb.UiSprite.EventHub.Publish(tasks.AllTasksCompleted{})
	}

	if wb.DstImg.Shader != nil {
		//turn off erase shader if it is completed
		maxCounter, ok := wb.DstImg.ShaderParams["MaxCounter"].(int)
		if ok {
			counter := wb.DstImg.ShaderParams["Counter"].(int)
			if counter >= maxCounter {
				wb.DstImg.Shader = nil
			}
		}
	}
	wb.UpdateTimers()
}

func (wb *WhiteBoardSprite) SubscribeOld(hub *tasks.EventHub) {
	hub.Subscribe(tasks.TaskRequirementsCompleted{}, func(e tasks.Event) {
		fmt.Printf("task requirment event recieved")
		if wb.wbState != WritingTask && wb.wbState != NotReceivingTasks {
			//wait to make cross out-able until writing is finished
			AddClickme(wb)
		} else {
			wb.clickAfterWriting = true
			//check for when writing is done and task can be crossed out
		}
		wb.taskJustCompleted = true
	})

	hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {

		ev := e.(tasks.TaskCreated)
		log.Printf("new task recieved @ whiteboard: %s", ev.Task.Text)
		wb.currentTask = ev.Task
		if wb.tasksCompletedToday == 0 {
			wb.timers[FirstTaskTimer].TurnOn()
		} else {
			wb.wbState = WritingTask
			wb.appendTextToOpenSlot(ev.Task.Text)
		}
	})

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		yCoord := float32(wb.spacing * (wb.tasksCompletedToday + 1))
		wb.AddCrossOutGraphicEntity(yCoord-10, 0, nil, defaultCrossOutShaderParams(), defaultCrossOutDrawOptParams())
		wb.checkAllTasksCompleted()

	})

	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		wb.initErase("")
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		wb.initErase("secondBuffer")
		wb.dayState = NightTime
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		fmt.Printf("White board received %d tasks on day %d", ev.NTasks, ev.Day)

		wb.wbState = IdleWB
		wb.dayState = JustStarted

		wb.numberOfTasksToday = ev.NTasks
		wb.tasksCompletedToday = 0

		wb.day = ev.Day

		wb.NoEraseDst.Shader = nil
		wb.NoEraseDst.Img.Clear()

		if ev.Day == 1 {
			wb.taskState = NotReceivingTasks
		}

		if ev.Day != 1 {
			wb.WriteHeader()
		}
	})

	hub.Subscribe(WriteToWhiteBoard{}, func(event tasks.Event) {
		ev := event.(WriteToWhiteBoard)
		if ev.Later {
			wb.msgsToWriteLater = append(wb.msgsToWriteLater, ev)
			if ev.EventDriven != nil {
				hub.Subscribe(ev.EventDriven, func(e tasks.Event) {
					wb.timers[WriteAfterTimer].TurnOn()
				})
			}
		} else {
			wb.wbState = WritingTask
			wb.appendTextToBestSpot(ev)
		}
	})
	hub.Subscribe(events.WindowOpened{}, func(event tasks.Event) {
		wb.windowOpen = true
	})
	hub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
		wb.windowOpen = false
		ev := e.(events.WindowClosed)
		if ev.Window == string(GrandpasJournal) && wb.day == 1 && !wb.JournalOnFirstDay {
			wb.WriteHeader()
			wb.JournalOnFirstDay = true
		}
	})
}

func (wb *WhiteBoardSprite) UpdateTimers() {

	for key, timer := range wb.timers {
		state := timer.Update()
		switch key {

		case AllTasksTimer:
			//Write out in marker effect/animation is completed
			if state == util.Done {
				timer.TurnOff()
			}

		case EraseTimer:
			//wipe away animation is completed
			if state == util.Done {
				for _, id := range wb.tIds {
					RemoveEntity(id)
				}
				wb.DstImg.Img.Clear()
				timer.TurnOff()
				if wb.dayState == NightTime {
					wb.appendTextToBestSpot(WriteToWhiteBoard{Msg: "All Done =)", PreferredPosition: "upperCenter"})
					wb.timers[UnderLineTimer].TurnOn()
				}
			}
		case FirstTaskTimer:
			if state == util.Done {
				timer.TurnOff()
			}
		case UnderLineTimer:
			if state == util.Done {
				rect := image.Rect(60, 32, 200, 48)
				sParams := defaultCrossOutShaderParams()
				sParams["Speed"] = 5.0
				dParams := defaultCrossOutDrawOptParams()
				wb.AddCrossOutGraphicEntity(20, 70, &rect, sParams, dParams)
				timer.TurnOff()
			}
		case WriteAfterTimer:
			if state == util.Done {
				if len(wb.msgsToWriteLater) <= 0 {
					return
				}
				wb.appendTextToBestSpot(wb.msgsToWriteLater[0])
				if wb.msgsToWriteLater[0].EventToPublish != nil {
					wb.timers[PublishTimer].TurnOn()
					wb.pubQueue = append(wb.pubQueue, wb.msgsToWriteLater[0].EventToPublish)
				}

				if len(wb.msgsToWriteLater) > 0 {
					wb.msgsToWriteLater = wb.msgsToWriteLater[1:]
				}

				if len(wb.msgsToWriteLater) > 0 {
					if wb.msgsToWriteLater[0].EventDriven != nil {
						timer.TurnOff()
					}
				} else {
					timer.TurnOff()
				}

				if len(wb.msgsToWriteLater) <= 0 {
					timer.TurnOff()
				}

			}
		case PublishTimer:
			if state == util.Done {
				wb.UiSprite.EventHub.Publish(wb.pubQueue[0])
				timer.TurnOff()
			}
		}
	}
}
