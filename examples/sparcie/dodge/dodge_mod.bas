10 REM text version of dodge
15 KEY OFF
20 DEFINT A-Z
30 RANDOMIZE TIMER
40 rem OPTION BASE 1
50 DIM MAP(75,22)
51 DIM BX(20), BY(20), BD(20)
52 GOSUB 11270
60 REM f$ = "autogen.txt"
70 REM GOSUB 11360
80 REM GOSUB 10300
90 NEXTT! = TIMER + .1
100 SCREEN 0: COLOR 2,0:CLS
110 PRINT TAB(33); "<<< Dodge >>>"
120 PRINT "Created by Sparcie.
130 PRINT "This game was originally a graphical platform game I wrote as a"
140 PRINT "sucessor to Jump, my first graphical game in basic."
150 PRINT "Basically you ("; :COLOR 9,0 : PRINT CHR$(1); : COLOR 2,0
160 PRINT ") need to avoid various enemies in the level"
170 PRINT "by jumping, moving, and using your jetpack." : PRINT
180 COLOR 6,0 : PRINT CHR$(22);:COLOR 2:PRINT " Chocolate gives you 500 points"
190 COLOR 11,0: PRINT CHR$(235);:COLOR 2:PRINT " The marble bag is 250 points
200 COLOR 12:PRINT CHR$(3);:COLOR 2:PRINT " The heart is for health and 5 points"
210 COLOR 10: PRINT CHR$(19);:COLOR 2: PRINT " Jetpack fuel so you can fly"
220 PRINT : PRINT "Controls are fairly simple"
230 PRINT "Using the numpad (because I'm lazy :p)"
240 PRINT " 4 - Left 5 - Stop 6 - Right"
250 PRINT " Space - Jump/Jetpack"
260 PRINT
270 PRINT  "Press a key"
280 C$ = INPUT$(1)
290 COLOR 2,0 : CLS
300 PRINT : PRINT "Select a level"
310 GOSUB 11270
320 PRINT "1-5 : Built in levels"
330 PRINT "C   : Load a custom level from a text file"
340 C$ = INPUT$(1)
350 IF C$="C" OR C$="c" THEN 450
360 IF C$="1" OR C$="2" OR C$="3" OR C$="4" OR C$="5" THEN 380
370 GOTO 340
380 IF C$="1" THEN RESTORE
390 IF C$="2" THEN RESTORE 30115
400 IF C$="3" THEN RESTORE 30230
410 IF C$="4" THEN RESTORE 30345
420 IF C$="5" THEN RESTORE 30460
425 PRINT :PRINT "Decoding level..."
430 GOSUB 11270: GOSUB 11360
440 GOTO 1000
450 COLOR 2,0 : CLS
460 PRINT
470 rem FILES "*.txt"
480 INPUT "Enter the text file to read> ",F$
490 ON ERROR GOTO 530
500 GOSUB 11270:GOSUB 10000
510 ON ERROR GOTO 0
520 GOTO 1000
530 PRINT "That file didn't exist (or there was another error)"
540 RESUME 480
1000 CLS: GOSUB 11000
1010 GOSUB 10300
1020 NEXTT! = TIMER + .1
1030 C$ = INKEY$
1035 IF HEALTH <= 0 THEN GOTO 2000
1040 IF C$ = "4" THEN PDIR = -1
1050 IF C$ = "6" THEN PDIR = 1
1060 IF C$ = "5" THEN PDIR = 0
1070 I = PYX : C = PYY : GOSUB 10500
1080 IF C$ = " " AND RESULT=1 AND JUF>0 THEN JUF=JUF-1 : UPF = 5
1090 IF C$ = " " AND RESULT<>1 THEN UPF=5
1100 IF NEXTT!>TIMER THEN 1030
1110 NEXTT! = TIMER + .1
1115 IF IMMUNE>0 THEN IMMUNE=IMMUNE-1
1120 GOSUB 11000
1130 GOSUB 9000
1140 GOSUB 9420 : GOSUB 8150
1150 GOSUB 9600
1160 GOSUB 9740
1170 ET = ET + 1
1180 IF ET = 2 THEN ET = 0
1190 IF ET= 1 THEN GOTO 1030
1200 GOSUB 7000
1210 IF CSY>-1 THEN GOSUB 8000
1220 GOSUB 8300
1230 IF CSY>-1 THEN GOSUB 8560
1240 GOSUB 9200
1250 GOTO 1030
2000 CLS
2010 PRINT "Ouch! It seems you have died."
2020 DIM HI$(10), SC(10)
2030 FOR I = 1 TO 10
2040 HI$(I) = "Nobody"
2050 SC(I)=0
2060 NEXT I
2070 ON ERROR GOTO 2150
2080 OPEN "dodge.top" FOR INPUT AS #1
2090 FOR I=1 TO 10
2100 INPUT #1,SC(I): INPUT #1,HI$(I)
2110 NEXT I
2120 CLOSE #1
2140 GOTO 2160
2150 RESUME 2160
2160 ON ERROR GOTO 0:IF SCORE <=SC(10) THEN 2240
2170 PRINT : PRINT
2180 INPUT "You got a high score! please enter your name >", N$
2190 FOR I = 1 TO 10
2200 IF SCORE<SC(I) THEN 2230
2210 SWAP SCORE,SC(I)
2220 SWAP HI$(I),N$
2230 NEXT I
2240 CLS
2250 PRINT :PRINT :PRINT :PRINT :PRINT
2260 FOR I = 1 TO 10
2270 COLOR I,0
2280 PRINT TAB(10),HI$(I), TAB(35), SC(I)
2290 NEXT I
2300 OPEN "dodge.top" FOR OUTPUT AS #1
2310 FOR I= 1 TO 10
2320 PRINT #1,SC(I) :PRINT #1,HI$(I)
2330 NEXT I
2340 CLOSE #1
2350 COLOR 2,0 : ERASE HI$, SC
2360 PRINT
2370 PRINT "Would you like to play again?"
2380 C$ = INPUT$(1)
2400 IF C$="Y" OR C$="y" THEN GOTO 290
2410 IF C$="N" OR C$="n" THEN END
2420 GOTO 2380
7000 IF BSX=-1 AND BSY=-1 THEN RETURN
7010 BSC=BSC-1
7020 IF BSC>0 THEN 7050
7021 IF BC=20 THEN 7050
7022 BSC = INT(RND*15)+6
7030 BC=BC+1
7040 BX(BC)=BSX:BY(BC)=BSY:BD(BC)=INT(RND*2)
7050 IF BC=0 THEN RETURN
7060 FOR T = 1 TO BC
7070 I=BX(T):C=BY(T)
7080 GOSUB 10200
7090 GOSUB 10500
7100 IF RESULT=0 THEN 7180
7110 BY(T)=BY(T)+1
7120 C=BY(T)
7130 GOSUB 10500
7140 IF RESULT=0 THEN 7220
7150 BD(T)=BD(T)+1
7160 IF BD(T)=2 THEN BD(T)=0
7170 GOTO 7220
7180 IF BD(T)=0 THEN GOSUB 10600 ELSE GOSUB 10700
7190 IF RESULT=0 THEN 7150
7200 IF BD(T)=0 THEN BX(T)=BX(T)-1
7210 IF BD(T)=1 THEN BX(T)=BX(T)+1
7220 IF BX(T)>1 AND BX(T)<75 THEN 7260
7230 BX(T)=BX(BC):BY(T)=BY(BC):BD(T)=BD(BC):BC=BC-1
7240 T=T-1
7250 GOTO 7290
7260 LOCATE BY(T)+1,BX(T) :COLOR 7,0
7270 PRINT "o"
7275 I = BX(T) : C = BY(T) : GOSUB 11560
7280 NEXT T
7290 RETURN
8000 IF CAX>0 AND CAY>0 THEN 8030
8010 CAX = CSX
8020 CAY = CSY
8030 I=CAX:C=CAY
8040 GOSUB 10500
8050 GOSUB 10200
8060 IF RESULT=1 THEN CAY=CAY+1:C=C+1
8070 GOSUB 10600
8080 IF RESULT=1 THEN CAX=CAX-1: GOTO 8110
8090 CAX = CSX
8100 CAY = CSY
8110 COLOR 10,0
8120 LOCATE CAY+1,CAX
8130 PRINT CHR$(145)
8131 GOSUB 8560
8135 I = CAX: C= CAY: GOSUB 11560
8140 RETURN
8150 IF JPT>0 THEN JPT = JPT - 1
8151 IF JPX<>-1 AND JPY<>-1 THEN 8160
8152 GOSUB 10400
8153 JPX = I : JPY = C
8160 IF JPT > 0 THEN RETURN
8170 I = JPX: C = JPY
8180 GOSUB 10500
8190 IF RESULT = 0 THEN 8220
8200 I = JPX: C = JPY: GOSUB 10200
8210 JPY = JPY +1
8220 LOCATE JPY+1, JPX: COLOR 10,0: PRINT CHR$(19)
8230 IF JPX<>PYX OR JPY<>PYY THEN RETURN
8240 JUF = JUF + 5
8250 IF JUF>15 THEN JUF=15
8260 JPT = INT(RND*200)+1
8270 GOSUB 10400
8280 JPX = I : JPY = C
8290 RETURN
8300 IF JMX <> -1 AND JMY <> -1 THEN 8340
8310 GOSUB 10400
8320 JMX = I: JMY = C
8330 GOTO 8350
8340 I = JMX: C = JMY : GOSUB 10200
8350 IF JMU = 0 OR JMY=1 THEN 8390
8360 JMU = JMU - 1
8370 IF MAP(JMX,JMY-1)<>1 AND MAP(JMX,JMY-1)<5 THEN JMY=JMY-1 ELSE JMU=0
8380 GOTO 8420
8390 I = JMX: C = JMY: GOSUB 10500
8400 IF RESULT = 1 THEN JMY = JMY + 1
8410 IF RESULT = 0 AND JMY > PYY THEN JMU = 7
8420 I = JMX : C = JMY
8430 IF PYX<=JMX THEN 8470
8440 GOSUB 10700
8450 IF RESULT = 1 THEN JMX = JMX + 1
8460 GOTO 8500
8470 IF PYX>=JMX THEN 8500
8480 GOSUB 10600
8490 IF RESULT = 1 THEN JMX = JMX - 1
8500 I = JMX : C = JMY: GOSUB 10900
8510 JMX = JMX + RESULT
8520 LOCATE JMY+1,JMX
8530 COLOR 12,0
8540 PRINT CHR$(21)
8545 I=JMX : C=JMY: GOSUB 11560
8550 RETURN
8560 IF MSF>0 THEN 8620
8570 MST = MST - 1
8580 IF MST>0 THEN RETURN
8590 MST = INT(RND*50)
8600 MSX = CAX: MSY = CAY
8610 MSF = 6
8620 I = MSX:C = MSY: GOSUB 10200
8630 MSF = MSF -1: RESULT =0
8640 IF PYX>MSX THEN GOSUB 10700
8650 IF PYX<MSX THEN GOSUB 10600
8660 IF RESULT = 0 THEN 8681
8670 IF PYX<MSX THEN MSX=MSX-1
8680 IF PYX>MSX THEN MSX=MSX+1
8681 IF MSY = 1 THEN 8691
8690 IF PYY<MSY AND MAP(MSX,MSY-1)<>1 AND MAP(MSX,MSY-1)<5 THEN MSY=MSY-1
8691 IF MSY = 22 THEN 8710
8700 IF PYY>MSY AND MAP(MSX,MSY+1)<>1 AND MAP(MSX,MSY+1)<5 THEN MSY=MSY+1
8710 LOCATE MSY+1,MSX
8720 COLOR 14
8730 IF MSF>0 THEN PRINT "*" : I = MSX : C=MSY: GOSUB 11560
8740 RETURN
9000 I = PYX
9010 C = PYY
9020 GOSUB 10200
9030 IF PDIR = -1 THEN GOSUB 10600
9040 IF PDIR = 1 THEN GOSUB 10700
9050 IF RESULT=1 THEN PYX = PYX+PDIR
9051 IF MAP(PYX,PYY)=2 THEN PDIR=0
9052 I=PYX:C=PYY: GOSUB 10900
9053 PYX=PYX+RESULT
9060 IF UPF = 0 THEN 9100
9070 UPF = UPF-1
9080 IF PYY=1 THEN 9100
9090 IF MAP(PYX,PYY-1)<>1 AND MAP(PYX,PYY-1)<5 THEN PYY=PYY-1
9100 IF UPF>0 THEN 9130
9110 I=PYX:C=PYY: GOSUB 10500
9120 IF RESULT=1 THEN PYY=PYY+1
9130 LOCATE PYY+1,PYX
9140 COLOR 9,0
9150 PRINT CHR$(1);
9160 RETURN
9200 IF GBX>0 AND GBY>0 THEN 9230
9210 GOSUB 10400
9220 GBX=I:GBY=C
9230 I=GBX:C=GBY
9240 GOSUB 10200
9250 GOSUB 10500
9260 IF RESULT = 0 THEN 9290
9270 GBY=GBY+1
9280 GOTO 9370
9290 GOSUB 10800
9300 IF RESULT = 0 OR GBY<=PYY THEN 9330
9310 GBY = GBY-1
9320 GOTO 9370
9330 IF GBX<PYX THEN GOSUB 10700
9340 IF GBX>PYX THEN GOSUB 10600
9350 IF RESULT=1 AND GBX<PYX THEN GBX=GBX+1
9360 IF RESULT=1 AND GBX>PYX THEN GBX=GBX-1
9370 I=GBX:C=GBY
9372 GOSUB 10900
9373 GBX=GBX+RESULT
9380 LOCATE GBY+1,GBX
9390 COLOR 4
9400 PRINT CHR$(64)
9405 I = GBX: C=GBY: GOSUB 11560
9410 RETURN
9420 IF HX>0 AND HY>0 THEN 9480
9430 HT = HT -1
9440 IF HT>0 THEN RETURN
9450 GOSUB 10400
9460 HX = I:HY = C
9470 HT = INT(RND*300) + 100
9480 COLOR 12,0
9490 LOCATE HY+1,HX: PRINT CHR$(3)
9500 IF HX <> PYX OR HY <> PYY THEN 9530
9510 IF HEALTH = 10 THEN 9530
9520 HEALTH = HEALTH + 1 : HT = INT(RND*300)+100
9521 I=HX:C=HY:GOSUB 10200
9522 HX=-1: HY=-1
9523 SCORE=SCORE+5
9530 HT = HT -1
9540 IF HT > 0 THEN RETURN
9541 I=HX:C=HY
9550 HX = -1
9560 HY = -1
9570 HT = INT(RND*300)+100
9571 GOSUB 10200
9580 RETURN
9600 IF MBX>0 AND MBY>0 THEN 9660
9610 IF MBT>0 THEN MBT=MBT-1
9620 IF MBT>0 THEN RETURN
9630 GOSUB 10400
9640 MBX = I
9650 MBY = C
9660 I=MBX:C=MBY:GOSUB 10500
9670 IF RESULT = 0 THEN 9700
9680 GOSUB 10200
9690 MBY=MBY+1
9700 LOCATE MBY+1,MBX
9710 COLOR 11
9720 PRINT CHR$(235)
9721 IF MBX<>PYX OR MBY<>PYY THEN RETURN
9722 SCORE=SCORE+250
9723 I=MBX:C=MBY:GOSUB 10200
9724 MBX=-1:MBY=-1
9725 MBT = INT(RND*300) + 100
9730 RETURN
9740 IF CHX>0 AND CHY>0 THEN 9790
9750 CHT=CHT-1
9760 IF CHT>0 THEN RETURN
9770 GOSUB 10400
9780 CHX=I:CHY=C
9790 I=CHX:C=CHY
9800 GOSUB 10500
9810 IF RESULT=0 THEN 9840
9820 GOSUB 10200
9830 CHY=CHY+1
9840 IF CHX<>PYX OR CHY<>PYY THEN 9890
9850 SCORE=SCORE+500
9860 CHX=-1:CHY=-1
9870 CHT=INT(RND*100)+100
9880 RETURN
9890 LOCATE CHY+1,CHX
9900 COLOR 6,0
9910 PRINT CHR$(22)
9920 RETURN
10000 OPEN F$ FOR INPUT AS #1
10010 I=1: C=1
10020 C$ = INPUT$(1,#1)
10025 IF I>75 THEN 10090
10030 IF C$ = "O" THEN MAP(I,C)=1
10040 IF C$ = "#" THEN MAP(I,C)=2
10050 IF C$ = "^" THEN MAP(I,C)=3 : CSX = I : CSY = C
10051 IF C$="v" THEN MAP(I,C)=4 : BSX=I: BSY=C
10060 IF C$ = "<" THEN MAP(I,C)=5
10070 IF C$ = ">" THEN MAP(I,C)=6
10080 IF C$ = " " THEN MAP(I,C)=0
10081 I=I+1
10090 IF C$ = CHR$(13) THEN I=0:C=C+1
10100 IF C<23 AND NOT(EOF(1)) THEN 10020
10110 CLOSE #1
10120 RETURN
10200 LOCATE C+1,I
10210 IF MAP(I,C)=0 THEN COLOR 7,0:PRINT " ";
10220 IF MAP(I,C)=1 THEN COLOR 13,0:PRINT CHR$(177);
10230 IF MAP(I,C)=2 THEN COLOR 7,0:PRINT  CHR$(197);
10240 IF MAP(I,C)=3 THEN COLOR 3,0:PRINT CHR$(127);
10241 IF MAP(I,C)=4 THEN COLOR 3,0:PRINT "v"
10250 IF MAP(I,C)=5 THEN COLOR 0,13:PRINT CHR$(174);
10260 IF MAP(I,C)=6 THEN COLOR 0,13:PRINT CHR$(175);
10270 RETURN
10300 FOR I = 1 TO 75
10310 FOR C = 1 TO 22
10320 GOSUB 10200
10330 NEXT C
10340 NEXT I
10350 RETURN
10400 I = INT(RND*75)+1
10410 C = INT(RND*22)+1
10420 IF MAP(I,C)>0 THEN 10400
10430 RETURN
10500 RESULT = 0
10510 IF C=22 THEN RETURN
10520 IF MAP(I,C+1) = 1 THEN RETURN
10530 IF MAP(I,C+1) = 2 OR MAP(I,C)=2 THEN RETURN
10540 IF MAP(I,C+1) =5 OR MAP(I,C+1)=6 THEN RETURN
10550 RESULT=1
10560 RETURN
10600 RESULT=1
10610 IF I=1 THEN RESULT=0:RETURN
10620 IF MAP(I-1,C)=1 OR MAP(I-1,C)>4 THEN RESULT=0
10630 RETURN
10700 RESULT=1
10710 IF I=75 THEN RESULT=0:RETURN
10720 IF MAP(I+1,C)=1 OR MAP(I+1,C)>4 THEN RESULT=0
10730 RETURN
10800 RESULT=1
10810 IF C=1 THEN RESULT=0:RETURN
10820 IF MAP(I,C)<>2 THEN RESULT=0
10830 IF MAP(I,C-1)<>2 AND MAP(I,C-1)<>0 THEN RESULT=0
10840 RETURN
10900 RESULT=0
10910 IF C=22 THEN RETURN
10920 IF MAP(I,C+1)=5 THEN RESULT=-1
10921 IF I=1 AND RESULT=-1 THEN RESULT=0
10930 IF MAP(I,C+1)=6 THEN RESULT=1
10931 IF I=75 AND RESULT=1 THEN RESULT=0
10940 IF RESULT = 0 THEN RETURN
10941 T = RESULT
10942 IF T = 1 THEN GOSUB 10700
10943 IF T = -1 THEN GOSUB 10600
10944 IF RESULT = 0 THEN RETURN
10945 RESULT = T
10946 RETURN
11000 LOCATE 1,1
11010 FOR I= 1 TO HEALTH
11020 COLOR 12,0
11030 PRINT CHR$(3);
11040 NEXT I
11050 FOR I = 1 TO 10-HEALTH
11060 COLOR 8
11070 PRINT CHR$(3);
11080 NEXT I
11090 PRINT "  ";
11100 FOR I = 1 TO JUF
11110 COLOR 10,0
11120 PRINT CHR$(24);
11130 NEXT I
11140 FOR I = 1 TO 15-JUF
11150 COLOR 8,0
11160 PRINT CHR$(24);
11170 NEXT I
11180 PRINT "  ";
11190 COLOR 7,0
11200 PRINT "Score ";SCORE;" ";
11210 C$ = CHR$(15)
11220 IF PDIR=-1 THEN C$=CHR$(17)
11230 IF PDIR=1 THEN C$=CHR$(16)
11240 COLOR 9
11250 PRINT C$;
11260 RETURN
11270 PYX=1 : PYY=1
11275 HEALTH = 10 : JUF = 10
11280 GBX=-1 : GBY=-1
11290 HX=-1:HY=-1: HT = 25 + INT(RND*200)
11295 MBX=-1 : MBY=-1 : MBT = 25 + INT(RND*200)
11300 CHX=-1 : CHY=-1 : CHT = 25 + INT(RND*200)
11310 CAX=-1 : CAY=-1 : CSX=-1 : CSY=-1
11320 JPX=-1 : JPY=-1 : JPT = 25 + INT(RND*200)
11330 JMX=-1 : JMY=-1 : JMU = 0
11340 MST=INT(RND*100) : MSF = 0
11345 BSX = -1 : BSY= -1 : BSC = INT(RND*15)+5
11350 RETURN
11360 C = 1
11370 WHILE C<23
11380 LC = 1
11390 WHILE LC<76
11400 READ IC, IT
11420 FOR I = LC TO LC + IC -1
11430 MAP(I,C) = IT
11440 NEXT I
11450 LC = LC + IC
11460 WEND
11470 C = C + 1
11480 WEND
11490 FOR I = 1 TO 75
11500 FOR C = 1 TO 22
11510 IF MAP(I,C) = 3 THEN CSX=I:CSY=C
11520 IF MAP(I,C) = 4 THEN BSX=I:BSY=C
11530 NEXT C
11540 NEXT I
11550 RETURN
11560 IF IMMUNE>0 THEN RETURN
11570 IF (I=PYX) AND (C=PYY) THEN HEALTH = HEALTH-1 : IMMUNE = 5
11580 RETURN
30000 REM lvla.txt
30005 DATA 49, 0, 1, 4, 25, 0
30010 DATA 32, 0, 9, 1, 1, 2, 4, 1, 24, 0, 1, 1, 1, 2, 3, 1
30015 DATA 8, 0, 12, 5, 21, 0, 1, 2, 29, 0, 1, 2, 3, 0
30020 DATA 41, 0, 1, 2, 28, 0, 1, 1, 1, 2, 3, 1
30025 DATA 24, 0, 13, 1, 4, 0, 1, 2, 29, 0, 1, 2, 3, 0
30030 DATA 41, 0, 1, 2, 29, 0, 1, 2, 3, 0
30035 DATA 41, 0, 1, 2, 12, 1, 17, 0, 1, 2, 3, 0
30040 DATA 41, 0, 1, 2, 29, 0, 1, 2, 3, 0
30045 DATA 1, 0, 1, 1, 1, 2, 16, 1, 16, 6, 6, 0, 1, 2, 29, 0, 1, 2, 3, 0
30050 DATA 2, 0, 1, 2, 38, 0, 1, 2, 29, 0, 1, 2, 3, 0
30055 DATA 2, 0, 1, 2, 38, 0, 1, 2, 15, 0, 11, 1, 3, 0, 1, 2, 3, 0
30060 DATA 2, 0, 1, 2, 38, 0, 1, 2, 29, 0, 1, 2, 3, 0
30065 DATA 8, 1, 1, 2, 10, 1, 22, 0, 1, 2, 8, 0, 11, 1, 1, 2, 13, 1
30070 DATA 8, 0, 1, 2, 5, 0, 15, 1, 12, 0, 1, 2, 19, 0, 1, 2, 13, 0
30075 DATA 8, 0, 1, 2, 32, 0, 1, 2, 19, 0, 1, 2, 13, 0
30080 DATA 8, 0, 1, 2, 32, 0, 1, 2, 19, 0, 1, 2, 12, 0, 1, 3
30085 DATA 5, 0, 13, 1, 2, 0, 11, 1, 10, 0, 1, 2, 19, 0, 1, 2, 13, 1
30090 DATA 30, 0, 2, 1, 1, 2, 7, 1, 1, 0, 1, 2, 19, 0, 1, 2, 13, 0
30095 DATA 32, 0, 1, 2, 8, 0, 1, 2, 1, 0, 17, 1, 1, 0, 1, 2, 13, 0
30100 DATA 32, 0, 1, 2, 8, 0, 1, 2, 12, 0, 7, 1, 1, 2, 13, 1
30105 DATA 32, 0, 1, 2, 8, 0, 1, 2, 19, 1, 1, 2, 3, 1, 10, 0
30110 DATA 32, 0, 1, 2, 8, 0, 1, 2, 19, 0, 1, 2, 13, 0
30115 REM lvlb.txt
30120 DATA 74, 0, 1, 2
30125 DATA 12, 0, 1, 2, 10, 1, 51, 0, 1, 2
30130 DATA 12, 0, 1, 2, 9, 0, 14, 1, 38, 0, 1, 2
30135 DATA 1, 0, 11, 1, 1, 2, 48, 0, 14, 1
30140 DATA 12, 0, 1, 2, 28, 0, 22, 1, 12, 0
30145 DATA 12, 0, 1, 2, 15, 0, 5, 1, 1, 2, 10, 1, 31, 0
30150 DATA 12, 0, 1, 2, 14, 0, 1, 2, 5, 1, 1, 2, 10, 1, 30, 0, 1, 1
30155 DATA 12, 0, 1, 2, 14, 0, 1, 2, 5, 0, 1, 2, 41, 0
30160 DATA 12, 0, 1, 2, 14, 0, 1, 2, 5, 0, 1, 2, 41, 0
30165 DATA 12, 0, 1, 2, 14, 0, 1, 2, 5, 0, 1, 2, 41, 0
30170 DATA 12, 0, 1, 2, 14, 0, 1, 2, 10, 1, 8, 0, 11, 1, 1, 2, 2, 1, 7, 0, 8, 1
30175 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30180 DATA 1, 0, 11, 1, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30185 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30190 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30195 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30200 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30205 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30210 DATA 12, 0, 1, 2, 14, 0, 1, 2, 20, 0, 9, 1, 1, 2, 7, 1, 10, 0
30215 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30220 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 17, 0
30225 DATA 12, 0, 1, 2, 14, 0, 1, 2, 29, 0, 1, 2, 16, 0, 1, 3
30230 REM lvlc.txt
30235 DATA 3, 0, 1, 2, 7, 0, 1, 2, 2, 0, 19, 1, 9, 0, 1, 4, 14, 0, 1, 2, 17, 0
30240 DATA 3, 0, 1, 2, 7, 0, 1, 2, 45, 0, 1, 2, 17, 0
30245 DATA 3, 0, 1, 2, 7, 0, 1, 2, 45, 0, 1, 2, 17, 0
30250 DATA 3, 0, 1, 2, 7, 0, 1, 2, 24, 0, 11, 1, 1, 2, 4, 1, 5, 0, 1, 2, 17, 0
30255 DATA 3, 0, 1, 2, 7, 0, 1, 2, 35, 0, 1, 2, 9, 0, 1, 2, 15, 0, 2, 1
30260 DATA 3, 0, 1, 2, 6, 0, 23, 1, 14, 0, 1, 2, 9, 0, 1, 2, 17, 0
30265 DATA 3, 0, 1, 2, 43, 0, 1, 2, 9, 0, 1, 2, 17, 0
30270 DATA 3, 0, 1, 2, 43, 0, 1, 2, 9, 0, 1, 2, 17, 0
30275 DATA 3, 0, 1, 2, 43, 0, 1, 2, 9, 0, 1, 2, 17, 0
30280 DATA 3, 0, 1, 2, 2, 0, 17, 1, 24, 0, 1, 2, 9, 0, 1, 2, 17, 0
30285 DATA 3, 0, 1, 2, 24, 0, 11, 5, 8, 0, 1, 2, 9, 0, 1, 2, 17, 0
30290 DATA 3, 0, 1, 2, 8, 0, 3, 6, 1, 2, 6, 1, 25, 0, 1, 2, 9, 0, 1, 2, 17, 0
30295 DATA 3, 0, 1, 2, 11, 0, 1, 2, 5, 0, 10, 1, 1, 2, 6, 1, 9, 0, 1, 2, 9, 0, 1, 2, 17, 0
30300 DATA 3, 0, 1, 2, 11, 0, 1, 2, 15, 0, 1, 2, 5, 0, 3, 1, 7, 0, 1, 2, 9, 0, 1, 2, 10, 0, 7, 1
30305 DATA 3, 0, 1, 2, 11, 0, 1, 2, 15, 0, 1, 2, 15, 0, 1, 2, 9, 0, 1, 2, 7, 0, 4, 1, 6, 0
30310 DATA 3, 0, 1, 2, 11, 0, 1, 2, 15, 0, 1, 2, 15, 0, 1, 2, 9, 0, 1, 2, 17, 0
30315 DATA 3, 0, 1, 2, 11, 0, 1, 2, 15, 0, 1, 2, 15, 0, 1, 2, 4, 0, 11, 1, 12, 0
30320 DATA 3, 0, 1, 2, 11, 0, 1, 2, 15, 0, 1, 2, 15, 0, 1, 2, 27, 0
30325 DATA 3, 0, 1, 2, 11, 0, 1, 2, 1, 0, 18, 1, 12, 0, 1, 2, 27, 0
30330 DATA 3, 0, 1, 2, 11, 0, 1, 2, 31, 0, 1, 2, 3, 0, 15, 1, 9, 0
30335 DATA 3, 0, 1, 2, 11, 0, 1, 2, 8, 0, 11, 1, 12, 0, 1, 2, 27, 0
30340 DATA 3, 0, 1, 2, 11, 0, 1, 2, 31, 0, 1, 2, 27, 0
30345 REM lvld.txt
30350 DATA 9, 0, 1, 2, 22, 0, 1, 2, 6, 0, 1, 4, 34, 0, 1, 3
30355 DATA 9, 0, 1, 2, 22, 0, 1, 2, 35, 0, 7, 1
30360 DATA 9, 0, 1, 2, 22, 0, 1, 2, 11, 0, 18, 1, 13, 0
30365 DATA 9, 0, 1, 2, 22, 0, 1, 2, 42, 0
30370 DATA 9, 0, 1, 2, 10, 0, 10, 1, 1, 2, 1, 0, 1, 2, 42, 0
30375 DATA 9, 0, 1, 2, 20, 0, 1, 2, 19, 1, 25, 0
30380 DATA 9, 0, 1, 2, 20, 0, 1, 2, 18, 0, 3, 1, 1, 2, 14, 1, 8, 0
30385 DATA 9, 0, 1, 2, 20, 0, 1, 2, 21, 0, 1, 2, 22, 0
30390 DATA 9, 0, 1, 2, 20, 0, 1, 2, 21, 0, 1, 2, 13, 0, 9, 1
30395 DATA 9, 0, 1, 2, 20, 0, 1, 2, 21, 0, 1, 2, 22, 0
30400 DATA 9, 0, 1, 2, 20, 0, 1, 2, 8, 0, 11, 1, 2, 0, 1, 2, 22, 0
30405 DATA 9, 0, 1, 2, 4, 0, 16, 1, 1, 2, 21, 0, 1, 2, 22, 0
30410 DATA 9, 0, 1, 2, 20, 0, 1, 2, 1, 0, 19, 1, 1, 0, 1, 2, 4, 0, 18, 1
30415 DATA 9, 0, 1, 2, 11, 0, 12, 1, 19, 0, 1, 2, 22, 0
30420 DATA 9, 0, 1, 2, 42, 0, 1, 2, 22, 0
30425 DATA 6, 0, 11, 1, 35, 0, 1, 2, 22, 0
30430 DATA 33, 0, 24, 1, 18, 0
30435 DATA 25, 0, 9, 1, 41, 0
30440 DATA 9, 0, 14, 1, 52, 0
30445 DATA 54, 0, 14, 1, 3, 0, 4, 1
30450 DATA 75, 0
30455 DATA 75, 0
30460 REM lvle.txt
30465 DATA 42, 0, 1, 4, 32, 0
30470 DATA 11, 0, 19, 1, 45, 0
30475 DATA 75, 0
30480 DATA 29, 0, 8, 1, 1, 2, 8, 1, 29, 0
30485 DATA 3, 0, 8, 1, 1, 2, 5, 1, 20, 0, 1, 2, 13, 0, 14, 1, 10, 0
30490 DATA 11, 0, 1, 2, 25, 0, 1, 2, 37, 0
30495 DATA 11, 0, 1, 2, 25, 0, 1, 2, 37, 0
30500 DATA 11, 0, 1, 2, 15, 0, 10, 1, 1, 2, 8, 1, 17, 0, 12, 1
30505 DATA 11, 0, 1, 2, 25, 0, 1, 2, 37, 0
30510 DATA 11, 0, 1, 2, 25, 0, 1, 2, 37, 0
30515 DATA 11, 0, 1, 2, 11, 0, 20, 1, 8, 0, 6, 1, 1, 2, 5, 1, 12, 0
30520 DATA 11, 0, 1, 2, 45, 0, 1, 2, 4, 0, 8, 1, 5, 0
30525 DATA 11, 0, 1, 2, 4, 0, 16, 1, 7, 0, 13, 1, 5, 0, 1, 2, 17, 0
30530 DATA 11, 0, 1, 2, 45, 0, 1, 2, 17, 0
30535 DATA 11, 0, 1, 2, 33, 0, 12, 1, 1, 2, 4, 1, 13, 0
30540 DATA 11, 0, 1, 2, 10, 0, 16, 1, 19, 0, 1, 2, 12, 0, 5, 1
30545 DATA 8, 0, 15, 1, 34, 0, 1, 2, 17, 0
30550 DATA 30, 0, 12, 1, 15, 0, 1, 2, 17, 0
30555 DATA 3, 0, 24, 1, 1, 2, 3, 1, 26, 0, 1, 2, 17, 0
30560 DATA 27, 0, 1, 2, 29, 0, 1, 2, 17, 0
30565 DATA 27, 0, 1, 2, 29, 0, 1, 2, 17, 0
30570 DATA 27, 0, 1, 2, 29, 0, 1, 2, 17, 0

