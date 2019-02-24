1 REM Digger Clone version 2 A Danson
2 REM    This is a clone of digger that I wrote as a kid sometime
3 REM    between 1990 and 1994. You'll notice the code is generally
4 REM    untidy, but hey I was only a kid.
5 REM Comments added today for read-ability. you can skip them if you like.
6 RANDOMIZE VAL(RIGHT$(TIME$,2))
7 screen 0:rem SCREEN 8
8 MBX=INT(RND*80)
9 MBY=1
10 CLS
20 KEY OFF
30 PRINT " digger"
40 PRINT
50 PRINT"written by andrew"
70 PRINT
80 INPUT "name";N$
81 LI=10:L=0
82 X=39
83 Y=9
90 CLS
100 COLOR 6
101 REM Displaying some base information for players
110 PRINT "information"
120 PRINT CHR$(178);" = dirt"
130 PRINT "@ = you
150 PRINT "# =enemys
160 PRINT "* =diamonds
170 PRINT "$ =money bag"
171 REM Simple top score mechanism
172 REM   replaces the highest score it can
173 REM   I still had a bit too learn about sorting etc at the time.
180 FOR I= 1 TO 10
190 IF SC>TSC(I) AND ST=0 THEN TSC(I)=SC:NA$(I)=N$:ST=1
200 LOCATE I,20
210 PRINT I;") ";NA$(I),TSC(I)
220 NEXT I
221 IF LI<=0 THEN SC=0:L=0:GOTO 80
222 ST=0
230 Z$=INPUT$(1)
231 IF J=0 THEN J=1 ELSE 280
232 REM Setup arrays for diamonds and enemies
260 dim ex(48),ey(48):rem DIM EX(16*3),EY(16*3)
270 dim dy(80),dx(80):rem DIM DY(16*5),DX(16*5)
271 REM Initialise the display
280 L=L+1
290 COLOR L
291 PLAY "l8mbabbcdefggfcdeafgedgfcd"
300 FOR XD= 1 TO 79
310 FOR YD=1 TO 20
320 LOCATE YD,XD
330 PRINT CHR$(178)
340 NEXT YD
350 NEXT XD
360 LOCATE 21
361 COLOR 9
370 PRINT "score ";SC,M$,N$
371 REM Create the diamonds and enemies for this level
380 ND=L*5
390 NE=L*3
400 FOR I=1 TO ND
410 LET DX(I)=INT(RND*79/2)*2+1
420 LET DY(I)=INT(RND*20/2)*2+1
430 LOCATE DY(I),DX(I)
440 COLOR L+2
450 PRINT "*"
460 NEXT I
470 FOR I=1 TO NE
480 EX(I)=INT(RND*79)+1
490 EY(I)=INT(RND*20)+1
500 LOCATE EY(I),EX(I)
510 COLOR L+3
520 PRINT "#"
530 NEXT I
531 REM not sure what these were meant to be for
540 ES=L+1
550 SP=L+2
551 IF L=12 THEN L=1:GOTO 280
552 LOCATE 22
553 PRINT "lives ";LI
555 REM Game loop
560 C$=INKEY$
570 IF C$="" AND T$=TIME$ THEN 560
580 M$=C$
590 LOCATE 22
600 PRINT "lives ";LI
610 LOCATE 21
620 PRINT "score ";SC,M$,N$
630 IF T$<>TIME$ THEN T=T+1:T$=TIME$:GOSUB 720
640 IF C$<>"" THEN GOSUB 660
641 IF LI=<0 THEN 90
650 GOTO 560
651 REM Keypress processing
652 REM   I often used the numeric keypad, cause I was lazy
653 REM   or early in the piece couldn't read arrow keys
660 IF C$="8" THEN D=1
670 IF C$="4" THEN D=2
680 IF C$="2" THEN D=3
690 IF C$="6" THEN D=4
700 IF C$="5" THEN D=0
710 RETURN
711 REM Enemy move
720 FOR I=1 TO NE
721 LOCATE EY(I),EX(I):PRINT " "
730 IF EX(I)>X THEN EX(I)=EX(I)-1
740 IF EX(I)<X THEN EX(I)=EX(I)+1
750 IF EY(I)>Y THEN EY(I)=EY(I)-1
760 IF EY(I)<Y THEN EY(I)=EY(I)+1
770 LOCATE EY(I),EX(I)
780 COLOR L+3
790 PRINT "#"
791 IF EX(I)=X AND EY(I)=Y THEN LI=LI-1:PLAY "mbb":GOSUB 2100
792 IF EX(I)=MBX AND EY(I)=MBY THEN PLAY "mbface":SC=SC+50:GOSUB 2100
800 NEXT I
801 REM Check end level, move the player, and draw diamonds
802 IF ND=0 THEN 280
803 FOR M=1 TO 2
804 LOCATE Y,X:PRINT " "
810 FOR I=1 TO ND
820 LOCATE DY(I),DX(I)
830 COLOR L+2
840 PRINT "*"
850 IF X=DX(I) AND Y=DY(I) THEN GOSUB 2000
860 NEXT I
870 IF D=1 THEN Y=Y-1
880 IF D=2 THEN X=X-1
890 IF D=3 THEN Y=Y+1
900 IF D=4 THEN X=X+1
901 IF X=>80 THEN X=1
902 IF X=<0 THEN X=79
903 IF Y=>21 THEN Y=1
904 IF Y=<0 THEN Y=20
910 LOCATE Y,X
920 COLOR L+1
930 PRINT "@"
931 NEXT M
932 REM Move the money bag down the screen
950 IF MBF=0 THEN MBF=INT(RND*2):RETURN
960 IF MBY=20 THEN MBF=0:MBX=INT(RND*80)+1:MBY=1:RETURN
962 LOCATE MBY,MBX:COLOR L:PRINT CHR$(177)
970 MBY=MBY+1
980 COLOR 14
990 LOCATE MBY,MBX
1000 PRINT "$"
1001 REM If you catch the money bag you get points!
1010 IF MBY=Y AND MBX=X THEN SC=SC+100:PLAY "mbdeadbeef"
1020 IF MBS=<37 THEN MBS=37+190
1021 MBS=MBS-10
1030 SOUND MBS,3
1040 RETURN
1991 REM Added so that the diamonds are removed correctly
1992 REM   Used as a subroutine seperately here as it's a later addition
1993 REM   sort of a side effect of hacking/working on something a while
2000 LOCATE DY(I),DX(I) : PRINT " "
2009 DX(I) = DX(ND)
2010 DY(I) = DY(ND)
2020 ND=ND-1
2030 SC=SC+10
2040 PLAY "mbcgcccgggg"
2050 RETURN
2091 REM Subroutine for removing enemies
2092 REM   Again could have been in the earlier enemy code
2093 REM   But hacked in later
2100 EY(I)=EY(NE):EX(I)=EX(NE):NE=NE-1
2110 RETURN
