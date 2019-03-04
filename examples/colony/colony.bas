1000 screen 0
1010 key off
1100 planetmax=20
1110 dim planetx(20), planety(20)
1120 banditmax=20
1130 dim banditx(20), bandity(20), banditd(20), banditt(20)
1200 planets=5
1210 bandits=5
1300 fieldsizex=80
1310 fieldsizey=23
1400 for i=1 to planets
1410 gosub 30000
1420 next
1500 for i=1 to bandits
1510 gosub 30100
1520 next
9000 h=1:hd=.1:h1=1
10000 rem main game loop
10010 playerx=int(fieldsizex / 2)
10020 playery=fieldsizey - 2
10030 while cmd$<>"q"
10040 cmd$=inkey$
10050 gosub 30200
10060 gosub 30300
10070 locate playery,playerx: print "T";
10080 gosub 30400
10090 gosub 30500
10900 wend
10999 end
30000 planetx(i) = int(rnd * fieldsizex) + 1
30010 planety(i) = int(rnd * fieldsizey) + 1
30020 return
30100 banditx(i) = int(rnd * fieldsizex) + 1
30110 bandity(i) = int(rnd * fieldsizey) + 1
30120 banditd(i) = i/2 : rem lower is faster
30190 return
30200 for i=1 to planets
30210 locate planety(i), planetx(i): print "O";
30220 next
30230 return
30300 for i=1 to bandits
30310 gosub 30800
30320 next
30330 return
30400 locate 24,20: print timer-t;"     ";
30403 if timer-t<hd then return
30405 t=timer
30408 m$=" hit q to exit "
30410 locate 25,h: print m$;
30430 if h<2 then h1=1
30440 if h+len(m$)>79 then h1=-1
30450 h=h+h1
30455 locate 24,10: print h;"  ";
30460 return
30500 for i=1 to bandits
30503 if timer-banditt(i)<banditd(i) then goto 30590
30505 banditt(i)=timer
30510 gosub 30600
30520 if planetx(p)<banditx(i) then gosub 30700:banditx(i)=banditx(i)-1:gosub 30800:goto 30590
30530 if planety(p)<bandity(i) then gosub 30700:bandity(i)=bandity(i)-1:gosub 30800:goto 30590
30540 if planetx(p)>banditx(i) then gosub 30700:banditx(i)=banditx(i)+1:gosub 30800:goto 30590
30550 if planety(p)>bandity(i) then gosub 30700:bandity(i)=bandity(i)+1:gosub 30800:goto 30590
30590 next
30599 return
30600 min=1000
30610 for j=1 to planets
30620 d=abs(planety(j)-bandity(i))+abs(planetx(j)-banditx(i))
30630 if d<min then min=d:p=j
30640 next
30650 return
30700 locate bandity(i), banditx(i): print " ";
30710 return
30800 locate bandity(i), banditx(i): print "#";
30810 return
