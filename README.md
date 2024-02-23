douyin
Douyin [Live Companion] push streaming key acquisition tool. Douyin live broadcast room barrage, room entry and other data acquisition, docking with Fay through Websocket

General idea
Obtain the rtmp address returned from broadcasting through a middleman agent

Implementation process
User installs CA certificate
Start proxy server
Broadcast detected
Parse the RTMP address
Force the end of the live broadcast companion (cannot click to disconnect)
OBS intervenes in streaming
Turn off proxy server
Exit this process
OpenSSL generates certificate
./certificates/generate-certificates.sh
MacOS Trust Certificate
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ./certificates/proxy-ca.crt
Windows Trust Certificate
Double-click to open ./certificates/proxy-ca.crt and click [Install Certificate]
Select [Local Computer] and click Next
Select [Put all certificates into the following storage (P)] -> Click [Browse]
Select the second [Trusted Root Certification Authority] -> click [OK] -> select [Next]
Select [Finish] -> Click [OK] -> Click [OK]
Android Trust Certificate
Send proxy-ca.crt to your mobile phone
. . . .
Ios trust certificate
Send proxy-ca.crt to your mobile phone
. . . .
Precautions
To turn off the live broadcast, please run the live broadcast companion again, click Continue live broadcast and then close the live broadcast. Otherwise, even if the stream is not pushed, it will not be downloaded immediately (I am too lazy to write down the broadcast)

In Fay, you need to change the uppercase first letter of json to lowercase
The code that needs to be modified in Fay

Operating Instructions
Docker run
Need to generate certificate first

./certificates/generate-certificates.sh
cd ./docker
docker-compose up -d
releaseDownload the executable file and run it
download
Generate certificate
trust certificate
Open executable file
Issue certificates to devices that require proxies
Set proxy address on computer/mobile phone ip:8001
grateful
goproxy
Fay
get-douyin-rtmp
