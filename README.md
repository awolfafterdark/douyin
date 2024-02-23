#douyin
Douyin [Live Companion] Push Streaming Key Obtaining Tool
Acquisition of data such as barrages in the Douyin live broadcast room, room entry, etc., connected to Fay through Websocket

### General idea
Obtain the rtmp address returned from broadcasting through a middleman agent

### Implementation process
1. User installs CA certificate
2. Start the proxy server
3. Broadcasting detected
4. Parse and get the RTMP address
5. Force the end of the live broadcast companion (cannot click to disconnect)
6. OBS intervenes in streaming
7. Turn off the proxy server
8. Exit this process


### OpenSSL generates certificate
```bash
./certificates/generate-certificates.sh
```

### MacOS Trust Certificate
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ./certificates/proxy-ca.crt
```
### Windows Trust Certificate
1. Double-click to open ./certificates/proxy-ca.crt and click [Install Certificate]
2. Select [Local Computer] and click Next
3. Select [Put all certificates into the following storage (P)] -> Click [Browse]
4. Select the second [Trusted Root Certification Authority] -> click [OK] -> select [Next]
5. Select [Finish] -> click [OK] -> click [OK]
###

### Android Trust Certificate
1. Send proxy-ca.crt to your mobile phone
2. . . . .

### Ios trust certificate
1. Send proxy-ca.crt to your mobile phone
2. . . . .

## Precautions
To turn off the live broadcast, please run the live broadcast companion again to continue the live broadcast and then close the live broadcast. Otherwise, even if the stream is not pushed, it will not be downloaded immediately (too lazy to write down the broadcast)

### In Fay, the capital letters of json need to be changed to lowercase.
![Code to be modified in Fay](fay/fay.pic.jpg)

## Run instructions

### Docker run
> Need to generate a certificate first
```bash
./certificates/generate-certificates.sh
cd ./docker
docker-compose up -d
```

### releaseDownload the executable file and run it

> 1. [Download](https://github.com/wwengg/douyin/releases)
> 2. Generate certificate
> 3. Trust certificate
> 4. Open the executable file
> 5. Send the certificate to the device that requires a proxy
> 6. Set proxy address ip:8001 on computer/mobile phone

### grateful
- [goproxy](https://github.com/elazarl/goproxy)
- [Fay](https://github.com/TheRamU/Fay)
- [get-douyin-rtmp](https://github.com/Cloud370/get-douyin-rtmp)
