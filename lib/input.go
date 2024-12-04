package lib 

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	gen "github.com/sizzlei/SimplePassGenerator"
	"errors"
)

type Passbox struct {
	Id 			int 
	Name 		string 
	Pass 		string
	Created 	string 
	Updated		string
}



func GetValue(m string) string {
	var v string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(m," : ")
	if scanner.Scan() {
		v = scanner.Text()	
	}
	
	return v
}

func GetValueInt(m string) (int, error) {
	var v int
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(m," : ")
	if scanner.Scan() {
		z := scanner.Text()	
		v, err = strconv.Atoi(z)
		if err != nil {
			return 0, err
		}
	}
	
	return v, nil
}


func NewConfigure() (Configure) {
	var cnf Configure 
	var err error
	cnf.Upper = GetValue("대문자 설정(기본값 : ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if cnf.Upper == "" {
		cnf.Upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	cnf.Lower = GetValue("소문자 설정(기본값 : abcdefghijklmnopqrstuvwxyz")
	if cnf.Lower == "" {
		cnf.Lower = "abcdefghijklmnopqrstuvwxyz"
	}

	cnf.UseDigits, err = GetValueInt("숫자 사용 여부(기본값 : 0(N))")
	if err != nil {
		cnf.UseDigits = 0
	}
	if cnf.UseDigits == 1 {
		cnf.Digits = GetValue("숫자 설정(기본값 : 0123456789")
		if cnf.Digits == "" {
			cnf.Digits = "0123456789"
		}
	}
	
	cnf.UseSpecial, err = GetValueInt("특수문자 사용 여부(기본값 : 0(N))")
	if err != nil {
		cnf.UseSpecial = 0
	}
	if cnf.UseSpecial == 1 {
		cnf.Special = GetValue("특수문자 설정(기본값 : !#$^&()><~")
		if cnf.Digits == "" {
			cnf.Digits = "!#$^&()><~"
		}
	}
	cnf.PassLength, err = GetValueInt("자동생성 패스워드 길이(기본값 : 12 / 최소 : 8)")
	if err != nil {
		cnf.UseSpecial = 12
	}

	if cnf.PassLength < 8{
		cnf.PassLength = 12
	} 

	return cnf
}

func Generate(cnf Configure) (string) {
	setSpec, setDigit := false,false
	if cnf.UseDigits == 1 {
		setDigit = true
	}
	if cnf.UseSpecial == 1 {
		setSpec = true
	}
	genertor := gen.New(
		gen.PassInOut{
			Upper:		cnf.Upper,
			Lower:		cnf.Lower,
			Digits:		cnf.Digits,
			Special:	cnf.Special,
			SpecialSet:	setSpec,
			DigitsSet: 	setDigit,
		},
	)

	pass, _ := genertor.GeneratePass(cnf.PassLength)
	
	return *pass
}

func NewPassBoxInput(cnf Configure) (Passbox, error) {
	var box Passbox 
	box.Name = GetValue("Box(Service) Name")
	if box.Name == "" {
		return box, errors.New("empty name")
	}
	box.Pass = NewPassInput(cnf)

	return box, nil
}

func NewPassInput(cnf Configure) string {
	var pass string
	genYn := GetValue("패스워드 자동 생성 여부(Y/N)")
	switch strings.ToLower(genYn) {
	case "y", "":
		pass = Generate(cnf)
	case "n":
		pass = GetValue("패스워드 입력")
	}

	return pass
}