package lib 

import (
	"os"
	"fmt"
)

const ( 
	appName = "PassBox"
	version = "1.0"
)

func Printlogo() {
	fmt.Println("+================================================+")
	fmt.Println("|")
	fmt.Println("| ",fmt.Sprintf("%s v%s",appName,version))
	fmt.Println("|")
	fmt.Println("|","by KURLY-ServiceEngineering")
	fmt.Println("|")
	fmt.Println("+================================================+")
	fmt.Println("")
	fmt.Println("")
}

func FileExsits(file string) bool {
	_, err := os.Stat(fmt.Sprintf("%s",file))
	if os.IsNotExist(err) {
		return false // 파일이 존재하지 않음
	}
	return true // 에러가 없으면 파일 존재
}

func FixView() {
	fmt.Print("\033[2J")    // 화면 지우기
	fmt.Print("\033[H")
}