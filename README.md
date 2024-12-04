# PassBox
### Description
Password Simple Storage
개인 로컬에 패스워드정보를 관리하기위해 간단하게 개발된 Password Storage 툴입니다. 

## Usage 
```shell
./passbox -storage=./passbox

+================================================+
|
|  PassBox v1.0
|
| by KURLY-ServiceEngineering
|
+================================================+


0 )  Create PassBox

1 )  PassBox View

2 )  Change Password configure

Enter Menu Idx :
```
### Initialize
최초 실행시 한번 진행되며, SQLITE DB 파일을 생성하고 설정합니다. 
```shell
INFO[0000] Start PassBox                                
INFO[0000] Passbox Storage : ./passbox
INFO[0002] initialize PassBox Storage                   
WARN[0002] 기본 설정값이 저장되어 있지 않습니다.                        
WARN[0002] 기본 설정을 시작합니다.                                
WARN[0002] 패스워드 생성시 필요한 정보를 설정 합니다.                     
대문자 설정(기본값 : ABCDEFGHIJKLMNOPQRSTUVWXYZ : 
소문자 설정(기본값 : abcdefghijklmnopqrstuvwxyz : 
숫자 사용 여부(기본값 : 0(N)) : 1
숫자 설정(기본값 : 0123456789 :  
특수문자 사용 여부(기본값 : 0(N)) : 1
특수문자 설정(기본값 : !#$^&()><~ :
```
### 신규 생성
> Passbox는 계정명을 저장하지 않습니다. 
```shell
0 )  Create PassBox

1 )  PassBox View

2 )  Change Password configure

Enter Menu Idx : 0
INFO[0111] 신규 패스워드 정보를 생성 합니다.                          
Box(Service) Name : Service1
패스워드 자동 생성 여부(Y/N) : y
자동생성 패스워드 길이(기본값 : 12 / 최소 : 8) : 12
INFO[0159] Added Box : Service1 [   Sly4XgTJ<6w4   ]
```
### 패스워드 보기
저장된 패스워드를 서비스 이름 별로 확인할 수 있습니다. 
```shell
0 )  Create PassBox

1 )  PassBox View

2 )  Change Password configure

Enter Menu Idx : 1

INFO[0051] Passbox [Service1] View                      
INFO[0051] Password : [  xPXvPuZfy70  ]                
INFO[0051] Last Updated : 2024-12-04 10:00:40           

0 )  Modify

1 )  Delete

2 )  History View

B ) Back
Enter Idx : 
```
### 저장된 패스워드 수정
패스워드를 변경해야할 경우 자동 완성 또는 수동으로 변경할 수 있습니다. 
```shell
INFO[0021] Passbox [Service1] View                      
INFO[0021] Password : [  xPXvPuZfy70  ]                
INFO[0021] Last Updated : 2024-12-04 10:00:40           

0 )  Modify

1 )  Delete

2 )  History View

B ) Back
Enter Idx : 0
패스워드 정보를 수정하시겠습니까?(Y/N) : 
패스워드 자동 생성 여부(Y/N) : n
패스워드 입력 : Test password
INFO[0035] 패스워드 변경이 완료되었습니다. [   Test password   ]      
Enter ) 돌아가기 : 
```

### 패스워드 변경 이력 조회 
패스워드를 변경하면 변경된 이력을 조회 할 수 있습니다. 
```shell
0 )  Modify

1 )  Delete

2 )  History View

B ) Back
Enter Idx : 2
INFO[0104] [ Service1 ] 변경 이력을 조회 합니다.                  
INFO[0104] Updated: 2024-12-04 10:08:58 Password: [  Test password  ] 
INFO[0104] Updated: 2024-12-04 10:00:40 Password: [  xPXvPuZfy70  ] 
Enter ) 돌아가기 : 
```

## 빠른사용 
빠르게 사용하기 위해서는 Alias 설정 후 단축어로 사용할 수 있습니다. (Mac 기준입니다.)
```
% ls -al | grep .zshrc
-rw-r--r--@  1 staff  staff         392 12  4 10:21 .zshrc

% vi .zshrc
alias passbox="/Password/PassBox_darwin_arm64 -storage=/Password/pass_storage


% source .zshrc
% passbox
```

또는 Sh 파일을 설정하여 실행합니다. 
./passbox.sh
```bash
#!/bin/bash

basedir="/Users/vnyx6xr63h/Desktop/Kurly-DBA/Password/PassBox_darwin_arm64"
storageFile="/Users/vnyx6xr63h/Desktop/Kurly-DBA/Password/pass_storage"

$basedir -storage=${storageFile}
```