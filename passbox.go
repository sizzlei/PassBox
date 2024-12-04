package main

import (
	"PassBox/lib"
	log "github.com/sirupsen/logrus"
	"fmt"
	// "github.com/rivo/tview"
	"strings"
	"time"
	"strconv"
	"flag"
)

func main() {
	var storage string 
	flag.StringVar(&storage,"storage","./passbox_storage","Passbox Storage")
	flag.Parse()

	lib.Printlogo()
	log.Info("Start PassBox")
	log.Infof("Passbox Storage : %s",storage)
	time.Sleep(2 * time.Second)



	// Storage 초기화
	fileBool := lib.FileExsits(storage)
	if fileBool == false {
		log.Info("initialize PassBox Storage")
		_, err := lib.InitDB(storage)
		if err != nil{
			log.Panic(err)
		}
	}

	// New Connection
	dbo, err := lib.ConnectStorage(storage)
	if err != nil {
		log.Panic(err)
	}
	
	// Init Setup
	cnf, err := dbo.CheckConfigure()
	if err != nil {
		if fmt.Sprint(err) == "sql: no rows in result set" {
			// 설정 정보가 없는 경우 신규 생성
			log.Warning("기본 설정값이 저장되어 있지 않습니다.")
			log.Warning("기본 설정을 시작합니다.")
			log.Warning("패스워드 생성시 필요한 정보를 설정 합니다.")
			cnf = lib.NewConfigure()

			err := dbo.SetConfigure(cnf)
			if err != nil {
				log.Panic(err)
			}
			
		} else {
			log.Panic(err)
		}
	} 

	// 기본 메뉴 정보
	menus := []string {
		"Create PassBox",
		"PassBox View",
		"Change Password configure",
	}

	passMenu := []string {
		"Modify",
		"Delete",
		"History View",
	}
	
	for {
		// 메인 메뉴
		MAINMENU:
		lib.FixView()
		lib.Printlogo()
		for i, v := range menus {
			fmt.Println(i,") ",v)
			fmt.Println("")
		}
		menuIdx := lib.GetValue("Enter Menu Idx")
		switch menuIdx {
		case "0":
			// ADD PASSWORD
			log.Info("신규 패스워드 정보를 생성 합니다.")
			box, err := lib.NewPassBoxInput(cnf)
			if err != nil {
				log.Error(err)
				time.Sleep(2 * time.Second)
				break
			}

			err = dbo.AddPassBox(box)
			if err != nil {
				log.Error(err)
			} else {
				_ = dbo.WriteHist(box)
				log.Infof("Added Box : %s [   %s   ]",box.Name, box.Pass)
				_ = lib.GetValue("Enter ) 돌아가기")
			}
		case "1":
			// VIEW PASSWORD
			log.Info("패스워드 정보를 확인 합니다.")

			PASSLIST:
			// Password 목록 조회
			box, err := dbo.GetPassList()
			if err != nil {
				log.Error(err)
			}

			if len(box) > 0 {
				// Password 목록
				lib.FixView()
				lib.Printlogo()
				for i, v := range box {
					fmt.Println(i,") ",v.Name)
					fmt.Println("")
				}
				fmt.Println("B ) Back")

				// 패스워드 선택
				boxId := lib.GetValue("Passbox Number")
				if strings.ToLower(boxId) == "b" {
					goto MAINMENU
				}

				boxIdx, err := strconv.Atoi(boxId)
				if err != nil {
					log.Error("잘못된 입력입니다.")
					_ = lib.GetValue("Enter ) 돌아가기")

					goto PASSLIST
				}

				if len(box) < boxIdx+1 {
					log.Warning("잘못된 번호를 입력 하였습니다.")
					_ = lib.GetValue("Enter ) 돌아가기")
					goto PASSLIST
				} else {
					lib.FixView()
					log.Infof("Passbox [%s] View",box[boxIdx].Name)
					log.Infof("Password : [  %s  ]",box[boxIdx].Pass)
					log.Infof("Last Updated : %s",box[boxIdx].Updated)
				}

				SUBMENU:
				fmt.Println("")
				for i, v := range passMenu {
					fmt.Println(i,") ",v)
					fmt.Println("")
				}
				fmt.Println("B ) Back")
				passIdx := lib.GetValue("Enter Idx")
				switch strings.ToLower(passIdx) {
				case "0":
					modifyPass := lib.GetValue("패스워드 정보를 수정하시겠습니까?(Y/N)")
					switch strings.ToLower(modifyPass) {
					case "y", "":
						newPass := lib.NewPassInput(cnf)
						err := dbo.UpdatePass(box[boxIdx].Id,newPass)
						if err != nil {
							log.Error(err)
							break
						}
						_ = dbo.WriteHist(lib.Passbox{
							Name: box[boxIdx].Name,
							Pass: newPass,
						})
						log.Infof("패스워드 변경이 완료되었습니다. [   %s   ]",newPass)
						_ = lib.GetValue("Enter ) 돌아가기")
						goto PASSLIST
					case "n":
						goto PASSLIST
					default:
						log.Error("잘못된 입력입니다.")
						_ = lib.GetValue("Enter ) 돌아가기")
						goto PASSLIST
					}
				case "1":
					deletePass := lib.GetValue("패스워드 정보를 삭제하시겠습니까?(Y/N)")
					switch strings.ToLower(deletePass) {
					case "y":
						err := dbo.DeletePass(box[boxIdx].Id)
						if err != nil {
							log.Error(err)
							break
						}
						log.Info("패스워드 삭제가 완료되었습니다.")
						_ = lib.GetValue("Enter ) 돌아가기")
						goto PASSLIST
					case "n":
						goto SUBMENU
					default:
						goto SUBMENU
					}
				case "2":
					log.Infof("[ %s ] 변경 이력을 조회 합니다.",box[boxIdx].Name)
					boxes, err := dbo.GetPassHist(box[boxIdx].Name)
					if err != nil{
						log.Error(err)
					}

					for _, v := range boxes{
						log.Infof("Updated: %s Password: [  %s  ]",v.Created,v.Pass)
					}
					_ = lib.GetValue("Enter ) 돌아가기")
					goto PASSLIST
				case "b":
					goto PASSLIST
				default:
					goto PASSLIST
				}

			} else {
				log.Warning("등록된 패스워드 정보가 없습니다. ")
				_ = lib.GetValue("Enter ) 돌아가기")
			}
		case "2":
			log.Info("패스워드 생성 설정 정보를 변경합니다.")
			cnf = lib.NewConfigure()
			
			err := dbo.SetConfigure(cnf)
			if err != nil {
				log.Panic(err)
			}

			log.Info("패스워드 생성 설정 정보가 재설정 되었습니다. ")
			_ = lib.GetValue("Enter ) 돌아가기")
		default:
			goto MAINMENU
		}
	}
	
}