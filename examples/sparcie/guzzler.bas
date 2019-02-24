1 DEFINT A-Z
2 RANDOMIZE INT(TIMER)
3 KEY OFF
10 REM ***********************
20 REM REMAKE OF GHOST GUZZLER
30 REM ***********************
40 REM I originally got this from a book
50 REM which I have since forgotten.
60 REM I've remimplemented it for better timing
70 REM and display.
77 screen 0
80 COLOR 5,7
90 CLS
100 PRINT "<<Ghost Guzzler>>"
110 PRINT "remade by Sparcie based on code from a long forgotten book"
120 PRINT
130 L$ = "," : R$ = "." : REM controls for up and down change to suite you
140 M$ = " " : REM control to mark number.
150 PRINT "The game is fairly simple, you change your number by pressing"
160 PRINT "the "+L$+" and "+R$+" keys to match the ghost crossing the screen."
170 PRINT "Then you press the space bar to capture (or guzzle) it."
180 PRINT "The quicker you do this the more points you will score."
190 PRINT
200 PRINT "Please press a number between 1 and 5 to select difficulty level"
210 C$ = INPUT$(1)
220 IF C$ = "1" OR C$ = "2" OR C$="3" OR C$="4" OR C$="5" THEN 240
230 GOTO 210
240 LEVEL = VAL(C$)
250 SCORE = 100 * ((LEVEL-1)*3)
260 LIVES = 3
270 GAP! = .5 / LEVEL : REM The length of a game tick in seconds.
280 CLS
290 NEXTT! = TIMER + GAP
300 GHOST = INT(RND*11)
310 PLAYER = 0
320 C$ = INKEY$ :REM start game loop
330 IF C$=L$ THEN PLAYER=PLAYER-1
340 IF C$=R$ THEN PLAYER=PLAYER+1
350 IF PLAYER = -1 THEN PLAYER = 10
360 IF PLAYER = 11 THEN PLAYER = 0
370 IF NOT(C$=" " AND PLAYER=GHOST) THEN 410
375 OLDSCORE = SCORE
380 SCORE = SCORE + (30 - GHOSTX)
385 IF (OLDSCORE MOD 100) > (SCORE MOD 100) THEN GAP! = GAP! - .025
390 GHOSTX = 1
400 GHOST = INT(RND*11)
410 IF NOT(GHOSTX=31) THEN 450
420 LIVES = LIVES-1
430 GHOSTX = 1
440 GHOST = INT(RND*11)
450 REM check timer and move ghost if needed
460 IF NOT(TIMER>=NEXTT!) THEN 490
470 GHOSTX = GHOSTX+1
480 NEXTT! = TIMER + GAP!
490 LOCATE 1
500 PRINT "Lives ";
510 FOR I = 1 TO LIVES
520 PRINT CHR$(3)+" ";
530 NEXT I
540 PRINT "  Score ";STR$(SCORE);"     "
550 PRINT TAB(GHOSTX); STR$(GHOST); TAB(30); ":"; STR$(PLAYER); "  "
560 PRINT TAB(29); "     "
570 IF LIVES=0 THEN 590
580 GOTO 320
590 REM end of game!
600 PRINT "You died!"
610 PRINT "Your score was "; STR$(SCORE)
620 PRINT
630 PRINT "Do you want to play again? (y/n)"
640 C$=INPUT$(1)
650 IF C$="y" THEN 200
660 IF C$="n" THEN 680
670 GOTO 640
680 PRINT "Bye, thanks for playing..."
690 END
