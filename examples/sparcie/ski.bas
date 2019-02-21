10 REM *************
20 REM Remake of Ski
30 REM *************
40 REM Another game from a long forgotten book.
50 REM It is interesting because it takes advantage of text scrolling
60 REM to move the play area. I've re-written the code for better timing
70 REM mostly as the original relied on CPU speed.
80 RANDOMIZE TIMER
90 DEFINT A-Z
99 screen 0: rem force screen 0
100 KEY OFF
110 COLOR 12,15 : REM change this if you don't like a white background
120 CLS
130 PRINT "<<< Downhill Skiiing >>>"
140 PRINT "Remade by Sparcie"
150 L$=",": R$=".": REM the controls
160 PRINT
170 PRINT "This is a simple game where you are skiing down a slope and"
180 PRINT "need to keep on the ski slop so you don't fall off the mountain."
190 PRINT "It's quite similiar to earlt car racing games."
200 PRINT "press "+L$+" to go left and "+R$+" to go right"
210 PRINT
220 PRINT "Select game speed (1 2 3 4)"
230 C$ = INPUT$(1)
240 IF C$="1" OR C$="2" OR C$="3" OR C$="4" THEN 260
250 GOTO 230
260 LEVEL = VAL(C$)
270 GAP! = .3/LEVEL : REM sets the gap between game ticks
280 NEXTT! = TIMER
290 SCORERATE = 5 * LEVEL
300 SCORE = 0
310 PLAYERX = 40
320 WID = 20
330 TRACK = 30
340 C$ = INKEY$
350 IF C$ = L$ AND PLAYERX>1 THEN PLAYERX=PLAYERX-1
360 IF C$ = R$ AND PLAYERX<79 THEN PLAYERX=PLAYERX+1
370 R = INT(RND*5)
380 IF R = 0 AND TRACK > 1 THEN TRACK=TRACK-1
390 IF R = 1 AND WID>5 THEN WID=WID-1
400 IF R = 3 AND WID<20 THEN WID=WID+1
410 IF R = 4 AND TRACK<79-WID THEN TRACK=TRACK+1
420 PRINT TAB(TRACK); "#"; TAB(PLAYERX); "!"; TAB(TRACK+WID); "#"
430 SCORE = SCORE + SCORERATE
440 REM a blank line!
450 IF TIMER<NEXTT! THEN 440
460 NEXTT! = TIMER+GAP!
470 IF PLAYERX=<TRACK OR PLAYERX>=(TRACK+WID) THEN 490
480 GOTO 340
490 PRINT "Ouch! you hit the wall!"
500 PRINT
510 PRINT "You got a score of: ",SCORE
520 PRINT "Do you wish to play again? (y/n)"
530 C$ = INPUT$(1)
540 IF C$ = "y" THEN 220
550 IF C$ = "n" THEN 570
560 GOTO 530
570 PRINT  "bye!..."

