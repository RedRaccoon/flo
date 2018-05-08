# flo

Download videos from flograppling on your pc

### Why this project

I'm working by bike/train, but there is no wifi in there ...  I still wanted to watch the videos so I had to found a solution

I really like flograppling but I often have trouble with the network :

When the traffic is slow the player load a low quality video and you can barely see something

I also wanted a solution to force the player to show a resolution of 1280 * 720

## Prerequisites

- having a pro account en flograppling

- if you want to build it from source you need golang 
https://golang.org/doc/install

and just go build ( no depedencies )

### Start


right click in "chrome" and click "inspect" to open chrome console (shortcut CTRL-SHIFT-i)
select the tab "network"
now your tab network is open you can to the video you want to watch

your chrome will download all the file to load the video 

you have to find the "playlist.m3u8" file, it's the one that stock all the ts files depending on the resolution 

(sometimes there is not the playlist.m3u8 like in , but I always found him for the ibjjf events )

right click - copy link url

![Alt text](ressources/readme/playlist3mu8.jpg?raw=true "Title")

then in cmd execut the 

flo.exe {url_of_the_playlist}

![Alt text](ressources/readme/godownload.jpg?raw=true "Title")

it will create a tmp files and stock all ts files inside it
you can check all the routine beeing launched

once it's done you will have a "final.ts" file with your video
