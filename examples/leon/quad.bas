10 REM *** WRITTEN BY LEON BARADAT  FEBRUARY 1984 ***
20 CLS:PRINT:LINE INPUT "What is your Self-Destruct Code? ";S$:SP1(0)=7.5
30 D$(0)="Warp":D$(1)="Impulse":D$(2)="Phasers":SP(1)=8:SP1(1)=8:SP1(2)=3000
40 D$(3)="Photon Torpedoes":D$(4)="Long-Range Scanner":SP(2)=3000:SP(0)=7.5
50 D$(5)="Short-Range Scanner":D$(6)="Status Report":D$(10)="Shields"
60 D$(7)="Computer":D$(8)="Transporter":D$(9)="Self-Destruct":PRINT
70 PRINT "  __";TAB(24);"____":O$="   1 2 3 4 5 6 7 8"
80 PRINT "--------------,   --------------,":E$="MR. SCOTT:  "
90 PRINT "'-------------'   '--/-/--------'":F$="MR. SPOCK:  "
100 PRINT TAB(10);"--'--------//--,":O1$="     1   2   3   4   5   6   7   8"
110 PRINT TAB(10);"'--------------'":GOSUB 2470:MS=RND*.5
120 PRINT:PRINT "    The Enterprise - NCC 1701":PRINT
130 PRINT:INPUT "How many Starbases (64 max) ";B
140 IF B<1 OR B>64 THEN GOSUB 2470:B=INT(5+RND*15)
150 REM *** SET UP THE GOODIES ***
160 GOSUB 2470:Q=INT(1+RND*8):GOSUB 2470:Q1=INT(1+RND*8):FOR I=1 TO B
170 GOSUB 2470:A=INT(1+RND*8):GOSUB 2470:B7=INT(1+RND*8)
180 Q(A,B7)=Q(A,B7)+10:IF Q(A,B7)>10 THEN Q(A,B7)=Q(A,B7)-10:GOTO 170
190 NEXT:INPUT "How many Klingons (576 max) ";K
200 IF K>0 AND K<577 THEN FOR J=1 TO K ELSE GOSUB 2470:K=INT(20+RND*30):GOTO 200
210 GOSUB 2470:A=INT(1+RND*8):GOSUB 2470:B7=INT(1+RND*8)
220 IF Q(A,B7)>=900 THEN 210 ELSE Q(A,B7)=Q(A,B7)+100
230 NEXT:FOR I=1 TO 8:FOR J=1 TO 8:GOSUB 2470:Q(I,J)=Q(I,J)+INT(RND*10)
240 NEXT:NEXT:PRINT:PRINT "Enter your name when ready to accept":T=5
250 LINE INPUT "Command of the Enterprise: ";M$:N$=M$+".":GOSUB 1740:E=3000:SH=50:K2=K
260 SD=VAL(RIGHT$(DATE$,2)+LEFT$(DATE$,2))+(VAL(MID$(DATE$,4,2)))/100:SD1=SD
270 PRINT:PRINT "It is Stardate";STR$(SD);".":PRINT:PRINT "We have been assigned a mission to":PRINT "seek out and destroy a fleet of";K
280 PRINT "Klingon Battle Cruisers in this":PRINT "section of the Galaxy.  We will":PRINT "have";B;"Starbases with which to":PRINT "resupply the Enterprise."
290 PRINT:PRINT "You will start out in Quadrant ";:Q2=Q:Q3=Q1:GOSUB 2480:PRINT ".":GOTO 1190
300 REM *** COMMAND ***
310 PRINT:A$(S,S1)="E":FOR A=0 TO 10:IF D(A)=0 THEN 330 ELSE NEXT
320 PRINT "MR. SCOTT:  We're dead in space!"
330 FOR A=0 TO 1:IF D(A)>0 AND D(A)<=1 AND SP(A)=SP1(A) THEN GOSUB 2470:SP(A)=INT(1+RND*7):PRINT E$;"Auxiliary power engaged.":PRINT " Maximum speed is ";D$(A);STR$(SP(A));"."
340 G(Q,Q1)=Q(Q,Q1):IF B1=1 THEN A$(B2,B3)="B" ELSE B1=0:B2=0:B3=0
350 NEXT:IF D(2)>0 AND D(2)<=1 AND SP(2)=SP1(2) THEN GOSUB 2470:SP(2)=INT(100+RND*900+.5):PRINT E$;"I can give you";SP(2);"units of Phaser power, ";N$
360 IF E<=20 THEN PRINT E$;"Power levels are dead, sir.":GOTO 2160
370 IF E<150 THEN PRINT E$;"Energy Levels are Critical, ";M$;"!!!"
380 IF Q(Q,Q1)=0 THEN G(Q,Q1)=.1
390 IF K<1 THEN 2240 ELSE PRINT "What are your Orders, ";M$;:INPUT C:C=INT(C)
400 GOSUB 2470:EM=INT(RND*21):C$="Green":GOSUB 1890:IF Q(Q,Q1)>=100 THEN C$="RED":K1=INT(Q(Q,Q1)/100)
410 IF C<0 OR C>9 THEN 420 ELSE ON C+1 GOTO 440,570,720,890,1120,1180,1220,1300,1580,1680
420 PRINT:PRINT "The Commands are:":FOR A=0 TO 9:PRINT "   #";RIGHT$(STR$(A),1);":  ";D$(A):NEXT:GOTO 310
430 REM *** WARP DRIVE ***
440 IF D(0)>1 THEN PRINT "Warp not available, ";N$:GOTO 310
450 PRINT:INPUT "Do you want to stay in the Quadrant";Q$:IF Q$="y" THEN 590
460 INPUT "To what Quadrant";Q2,Q3:W3=SQR((Q2-Q)^2+(Q3-Q1)^2)
470 IF Q2<1 OR Q3<1 OR Q2>8 OR Q3>8 OR W3=0 THEN 310
480 PRINT:PRINT "Distance to Quadrant ";:GOSUB 2480:PRINT " -";W3
490 PRINT:PRINT D$(C);" Factor";:INPUT W1:IF W1<=0 THEN 310
500 IF C=1 THEN 510 ELSE IF W1>7 AND D(0)=0 THEN MS=(INT(MS*100+.5)/100):PRINT:PRINT E$;M$;", I think I can manage Warp";STR$(7.5+MS);".":W1=7.5+MS:SS=1:GOTO 520
510 IF W1>SP(C) THEN PRINT:PRINT "Maximum speed is ";D$(C);STR$(SP(C));"!":GOTO 490
520 W2=W3/W1^2:ED=W1*50:TA=TA+W2:IF C=1 THEN TA=TA-W2:W2=W3*10/W1:ED=W1*10:TA=TA+W2
530 IF W2>2 THEN PRINT" MR. SPOCK:  This will take over"INT(W2)"Stardates":PRINT" at "D$(C);STR$(W1)".  Proceed";:INPUT A$:IF LEFT$(A$,1)<>"Y" AND LEFT$(A$,1)<>"y" THEN 310
535 ED=INT(ED+.5):PRINT:PRINT "Course plotted.  Energy drain:";STR$(ED);".":E=E-ED-EM
540 PRINT "Estimated time to arrival:";W2;"Stardates.":IF SS=1 THEN GOSUB 2490:SS=0
550 FOR A=1 TO 7500*W2:NEXT:PRINT:PRINT "Now entering Quadrant ";:GOSUB 2480:PRINT " . . .":Q=Q2:Q1=Q3:ST=0:GOSUB 1740:GOTO 1930
560 REM *** IMPULSE ***
570 IF D(1)>1 THEN PRINT "Impulse not available, ";N$:GOTO 310
580 PRINT:INPUT "Do you want to leave the Quadrant";Q$:IF LEFT$(Q$,1)="y" THEN 460
590 INPUT "To what Sector";Q2,Q3:I2=SQR((Q2-S)^2+(Q3-S1)^2)
600 IF Q2<1 OR Q3<1 OR Q2>8 OR Q3>8 OR I2=0 THEN 310
610 PRINT:PRINT "Distance to Sector ";:GOSUB 2480:PRINT " -";I2
620 PRINT:PRINT D$(C);" Factor";:INPUT I1:IF I1<=0 THEN 310
630 IF C=1 THEN 640 ELSE IF I1>7 THEN PRINT:PRINT E$;M$;", I think I can manage Warp";STR$((7.5+MS)/10);".":I2=7.5+MS:GOTO 650
640 IF I1>SP(C) THEN PRINT:PRINT "Maximum speed is ";D$(C);STR$(SP(1));"!":GOTO 620
650 A$(S,S1)=" ":S=Q2:S1=Q3:IF A$(S,S1)="*" THEN 2260
660 IF A$(S,S1)="K" THEN A$(S,S1)="E":GOTO 2310
670 IF A$(S,S1)="B" THEN A$(S,S1)="E":GOTO 2410
680 PRINT:PRINT "Course plotted.  Energy drain:";STR$(I1*10);".":E=E-I1*10+EM:TA=TA+I2/I1
690 PRINT "Estimated time to arrival:";I2/I1;"Stardates.":FOR A=1 TO 7500*(I2/I1):NEXT
700 A$(S,S1)=" ":PRINT:PRINT "Now in Sector ";:GOSUB 2480:PRINT ".":A$(S,S1)="E":GOTO 1930
710 REM *** PHASERS ***
720 PRINT:IF D(2)>1 THEN PRINT "Phasers inoperative, ";N$:GOTO 310
725 IF D(7)>0 THEN PRINT "Phasers will not be at full efficacy..."
730 IF K1=0 THEN PRINT F$;"Sensors show no Klingons in this Quadrant.":GOTO 310
740 PRINT "Phasers locked on Target.  Energy available:";E;"units."
750 PRINT:INPUT "Number of units to fire";X:IF X<=0 OR X>E THEN 310
760 A=0:IF X>SP(2) THEN PRINT "All I can give you is";SP(2);"units, ";M$;"!":GOTO 750
770 TA=TA+.1:E=E-X-EM:IF D(7)>0 THEN GOSUB 2470:PP=INT(20+(RND*.6)*100):PRINT "MR. CHEKOV:  Phaser lock is malfunctioning.":PRINT TAB(14);"Phasers working at";STR$(PP);"% capacity.":X=X*(PP/100)
780 A=A+1:GOSUB 2470:H=INT(3*((1+RND*.2-.4)*X/((P(A)-S)^2+(N(A)-S1)^2)^.25+.5)/K1):FOR B7=1 TO 3:PRINT CHR$(7);:NEXT
790 KE(A)=KE(A)-H:PRINT H;"unit hit on Klingon at Sector ";:Q2=P(A):Q3=N(A):GOSUB 2480:PRINT "."
800 IF KE(A)<=0 THEN PRINT "    *** Klingon Destroyed ***":GOTO 850
810 GOSUB 2470:IF RND*50>(KE(A)/10) OR KH>4 THEN 830
820 GOSUB 2470:B7=INT(1+RND*3):GOSUB 2470:KD(A,B7)=KD(A,B7)+RND*4+.4:PRINT "Klingon's ";D$(B7);" damaged, ";N$:KH=KH+1:GOTO 810
830 KH=0:IF KD(A,1)>0 AND KD(A,2)>0 AND KD(A,3)>0 THEN PRINT F$;"Klingon disabled, ";N$
840 PRINT "   Sensors show";KE(A);"units remaining.":GOTO 860
850 K=K-1:A$(P(A),N(A))=" ":FOR AK=A TO K1:P(AK)=P(AK+1):N(AK)=N(AK+1):T(AK)=T(AK+1):T1(AK)=T1(AK+1):KE(AK)=KE(AK+1):NEXT:T3=T3-1:T2=T2-1:Q(Q,Q1)=Q(Q,Q1)-100:A$(P(K1),N(K1))=" ":K1=K1-1:A=A-1:GOSUB 2570
860 IF A<K1 THEN 780 ELSE IF X<1001 THEN 1930 ELSE GOSUB 2470:IF RND*(X/500)<1.5 THEN 870 ELSE 1930
870 PRINT E$;"Phasers have overloaded, ";N$:GOSUB 2470:D(2)=D(2)+RND*4:PRINT " Time to repair -";D(2):GOTO 310
880 REM *** PHOTON TORPEDO ***
890 IF D(3)>0 THEN PRINT D$(3);" are under repair, ";N$:GOTO 310
900 IF T<1 THEN PRINT "We're out of ";D$(3);", ";N$:GOTO 310
910 IF T2<1 THEN 1090 ELSE PRINT:PRINT D$(3);" locked on Target.  Ready to fire on Command."
920 LINE INPUT "Hit return to release Torpedoes.  ";T$:IF T$<>"" THEN 310
930 TA=TA+.1:E=E-EM-50:A=1
940 IF T<1 THEN 900 ELSE IF A$(T(A),T1(A))="B" THEN PRINT F$;"I fail to see the logic in destroying our Starbase, ";N$:B8=B8+1:Q(Q,Q1)=Q(Q,Q1)-10:A$(T(A),T1(A))=" ":T=T-1:B1=0:GOTO 1100
950 IF A$(T(A),T1(A))=" " THEN PRINT "MR. SULU:  There's nothing to fire at in Sector ";:Q2=T(A):Q3=T1(A):GOSUB 2480:PRINT ", ";N$:GOTO 1100
960 IF A$(T(A),T1(A))="E" THEN PRINT F$;M$;", ordering us to shoot ourselves is not logical.":GOTO 1100
970 GOSUB 2470:H=INT(250*(1+RND*.5)+.5):IF A$(T(A),T1(A))="*" THEN PRINT F$;"Destroying stars again, ";M$;"?":A$(T(A),T1(A))=" ":Q(Q,Q1)=Q(Q,Q1)-1:T=T-1:GOTO 1100
980 GOSUB 2470:B6=RND*2:IF B6<SQR((P(A)-S)^2+(N(A)-S1)^2)/(10+D(7)) THEN 1080
990 AK=A:FOR AR=1 TO K1:IF P(AR)=T(A) AND N(AR)=T1(A) THEN AK=AR
1000 NEXT
1010 IF A$(T(A),T1(A))<>"K" OR A>T3 THEN 1100
1020 T=T-1:PRINT CHR$(7);:KE(AK)=KE(AK)-H:PRINT H;"unit hit on Klingon at Sector ";:Q2=T(A):Q3=T1(A):GOSUB 2480:PRINT ".":IF KE(A)<=0 THEN PRINT "    *** Klingon Destroyed ***":T2=T2-1:GOTO 1070
1030 GOSUB 2470:IF RND*50>(KE(AK)/10) OR KH>4 THEN 1050
1040 GOSUB 2470:B7=INT(1+RND*3):GOSUB 2470:KD(AK,B7)=KD(AK,B7)+RND*4+.4:PRINT "Klingon's ";D$(B7);" damaged, ";N$:KH=KH+1:GOTO 1030
1050 IF KD(AK,1)>0 AND KD(AK,2)>0 AND KD(AK,3)>0 THEN PRINT F$;"Klingon disabled, ";N$
1060 PRINT "   Sensors show";KE(A);"units remaining.":GOTO 1100
1070 K=K-1:A$(T(A),T1(A))=" ":FOR AK=A TO K1:P(AK)=P(AK+1):N(AK)=N(AK+1):T(AK)=T(AK+1):T1(AK)=T1(AK+1):KE(AK)=KE(AK+1):NEXT:T3=T3-1:T2=T2-1:Q(Q,Q1)=Q(Q,Q1)-100:A$(P(K1),N(K1))=" ":K1=K1-1:A=A-1:GOSUB 2570:GOTO 1100
1080 O1=1:IF P(AK)<=0 AND N(AK)<=0 THEN 1010 ELSE PRINT F$;"Klingons at ";:Q2=P(AK):Q3=N(AK):GOSUB 2480:PRINT " have outmaneuvered our Torpedo.":T=T-1:GOSUB 2530:GOTO 1100
1090 PRINT D$(3);" are not locked, ";N$:INPUT "Where should we fire";T(1),T1(1):IF T(1)<=0 AND T1(1)<=0 THEN 310 ELSE T3=1:GOTO 920
1100 IF A>=T3 THEN 1930 ELSE A=A+1:GOTO 940
1110 REM *** LONG-RANGE SCANNER ***
1120 IF D(4)>0 THEN PRINT D$(4);" is under repair, ";N$:GOTO 310
1130 PRINT "     ";Q1-1;" ";Q1;" ";Q1+1:FOR A=Q-1 TO Q+1:PRINT "    -------------":PRINT A;" :";:FOR B7=Q1-1 TO Q1+1
1140 IF A<1 OR A>8 OR B7<1 OR B7>8 THEN PRINT "***";:GOTO 1160
1150 IF Q(A,B7)=.1 THEN PRINT "000"; ELSE PRINT RIGHT$(STR$(Q(A,B7)+1000),3);
1160 PRINT ":";:G(A,B7)=Q(A,B7):GOSUB 2550:NEXT:PRINT " ";A:NEXT:PRINT "    -------------":PRINT "     ";Q1-1;" ";Q1;" ";Q1+1:GOTO 310
1170 REM *** SHORT-RANGE SCAN ***
1180 IF D(5)>0 THEN PRINT D$(5);" is under repair, ";N$:GOTO 310
1190 PRINT:PRINT O$:TA=TA+.1:FOR C=1 TO 8:PRINT STR$(C);:FOR D=1 TO 8
1200 PRINT " ";A$(C,D);:NEXT:PRINT C:NEXT:PRINT O$:GOTO 310
1210 REM *** STATUS REPORT ***
1220 IF D(6)>0 THEN PRINT D$(6);" can't get through to us, ";N$:GOTO 310
1230 PRINT "Condition --------- ";C$:PRINT "Reserve Energy ----";E
1240 PRINT "Stardate ----------";SD:PRINT "Quadrant ----------";Q;",";Q1
1250 PRINT "Sector ------------";S;",";S1:PRINT "Stardates Passed --";SD-SD1
1260 PRINT "Shield Energy -----";STR$(SH);"%%":PRINT D$(3);" --";T
1270 PRINT "Klingons Left -----";K:PRINT "Starbases ---------";B
1280 PRINT:FOR A=0 TO 10:PRINT D(A),D$(A):NEXT:TA=TA+.1:GOTO 310
1290 REM *** COMPUTER ***
1300 PRINT:PRINT "Computer ready.  Which option do you wish, ";M$;:INPUT O:O=INT(O)
1310 IF D(7)>0 THEN IF O=2 OR O=3 OR O=4 OR O=6 OR O=7 THEN PRINT "The Computer is smashed for the time being, ";N$:GOTO 310
1320 IF O=3 AND D(10)>0 THEN PRINT "Shields inoperative, ";N$:GOTO 310
1330 IF O<1 OR O>8 THEN 1340 ELSE PRINT:ON O GOTO 1380,1390,1420,1450,1490,1510,1520,1560
1340 PRINT "   #1:  Change Name":PRINT "   #2:  Distance Calculator"
1350 PRINT "   #3:  Shield Energy Level":PRINT "   #4:  Probe"
1360 PRINT "   #5:  Rest":PRINT "   #6:  Lock Torpedoes"
1370 PRINT "   #7:  Galactic Record":PRINT "   #8:  Clear Torp Lock":PRINT "Which Option, ";M$;:INPUT O:O=INT(O):GOTO 1310
1380 PRINT "What do you want to change your name to, ";M$;:INPUT M$:N$=M$+".":GOTO 310
1390 PRINT "We are in Quadrant ";:Q2=Q:Q3=Q1:GOSUB 2480:PRINT TAB(13);"Sector ";:Q2=S:Q3=S1:GOSUB 2480
1400 PRINT:INPUT "  Enter Initial Coordinates:  ";V,V1:INPUT "  Enter Second Coordinates:  ";Q2,Q3
1410 PRINT "Distance to Location ";:GOSUB 2480:PRINT " - ";SQR((V-Q2)^2+(V1-Q3)^2):GOTO 310
1420 PRINT "Shields presently at";STR$(SH);"%  Enter new Shield Energy Level: ";:INPUT SN
1430 IF SN<0 OR SN>100 THEN PRINT "Sorry, but the Shields can't take that percentage.":GOTO 310
1440 SH=INT(SN):GOTO 310
1450 PRINT "MR. SULU:  Probe launched, ";N$:PRINT:PRINT "There are";ST;"stars and";B1;"Starbases in this Quadrant."
1460 IF K1=0 THEN PRINT "There are no Klingons in the Quadrant, ";N$:GOTO 310
1470 PRINT "Probe indicates: ";K1;"Klingons in this Quadrant.":PRINT:PRINT "SECTOR","ENERGY","TORPEDOES"," DISTANCE"
1480 FOR A=1 TO K1:PRINT " ";:Q2=P(A):Q3=N(A):GOSUB 2480:PRINT KE(A),"   ";KT(A),SQR((P(A)-S)^2+(N(A)-S1)^2):NEXT:GOTO 310
1490 INPUT "MR. SULU:  So just how long are we gonna sit out here in space";RP:IF RP<=0 OR RP>5 THEN PRINT "No jokes please, ";N$:GOTO 310
1500 FOR A=1 TO RP*7500:NEXT:TA=TA+RP:GOTO 1930
1510 T2=K1:T3=K1:FOR A=1 TO K1:T(A)=P(A):T1(A)=N(A):NEXT:PRINT "Photon Torpedoes locked on all targets.":GOTO 310
1520 PRINT " MR. SULU:  We are in Quadrant"Q","Q1".":PRINT
1530 PRINT O1$:FOR A=1 TO 8:PRINT A;:FOR B7=1 TO 8:PRINT " ";:IF G(A,B7)=.1 OR Q(A,B7)=.1 THEN PRINT "000";:GOTO 1550
1540 IF G(A,B7)=0 THEN PRINT "***"; ELSE PRINT RIGHT$(STR$(G(A,B7)+1000),3);
1550 NEXT:PRINT A:NEXT:PRINT O1$:GOTO 310
1560 T3=0:T2=0:PRINT F$;"Torpedo Lock cleared, ";N$:GOTO 310
1570 REM *** DOCK  W / STARBASE ***
1580 IF D(8)>0 THEN PRINT D$(8);" is damaged.  Sorry, ";N$:GOTO 310
1590 IF SQR((B2-S)^2+(B3-S1)^2)>2 OR B1=0 THEN PRINT "Sorry, ";M$;", but you can't Dock without a Starbase!":GOTO 310
1600 PRINT:PRINT F$;"We have Docked.  All energy levels back to normal.":TA=TA+.1
1610 DD=0:FOR A=0 TO 10:DD=DD+D(A):NEXT:IF DD<=0 THEN 1660
1620 PRINT E$;"Starfleet Technicians are standing by to beam up":PRINT TAB(13);"and effect repairs on the following:"
1630 FOR A=0 TO 10:IF D(A)<>0 THEN PRINT TAB(15);D$(A)
1640 NEXT:DD=DD/10:PRINT TAB(12);DD;"Stardates estimated repair time.":PRINT TAB(13);"Will you authorize the repair orders, "M$;:INPUT A$
1650 IF LEFT$(A$,1)="y" OR LEFT$(A$,1)="Y" THEN VA=DD*(RND*.4+.8499999):TA=TA+VA:FOR A=1 TO VA*7500:NEXT:FOR A=0 TO 10:D(A)=0:NEXT:PRINT:PRINT "Repairs completed, "N$
1660 E=3000:T=5:SP(0)=9:SP(1)=8:SP(2)=3000:GOTO 1980
1670 REM *** SELF-DESTRUCT ***
1680 IF D(9)>0 THEN PRINT D$(9);" Inoperative, ";N$:GOTO 310
1690 PRINT:PRINT D$(9);" ready, ";N$;"  Enter ";D$(9);" Code: ";:LINE INPUT S1$
1700 IF S1$<>S$ THEN GOSUB 2470:PRINT:PRINT "Incorrect ";D$(9);" Code.":D(9)=D(9)+RND*4:GOTO 310
1710 E=E-EM:PRINT:FOR A=10 TO 1 STEP -1:PRINT TAB(19-INT(A/10));A;CHR$(7):FOR B7=1 TO 650:NEXT B7:NEXT A:K=K-K1:B8=B8+B1
1720 FOR A=1 TO 4:PRINT CHR$(7):NEXT:PRINT "BOOM!!!  You have been splattered all over the Galaxy!":GOTO 2160
1730 REM *** ENTERING QUADRANT ***
1740 O1=0:FOR D=1 TO 8:FOR C=1 TO 8:A$(D,C)=" ":NEXT:NEXT
1750 ST=VAL(RIGHT$(STR$(Q(Q,Q1)),1)):FOR A=1 TO ST
1760 GOSUB 2470:I=INT(1+RND*8):GOSUB 2470:J=INT(1+RND*8)
1770 A$(I,J)="*":NEXT
1780 GOSUB 2470:S=INT(1+RND*8):GOSUB 2470:S1=INT(1+RND*8)
1790 IF A$(S,S1)="*" THEN 1780
1800 A$(S,S1)="E":K1=INT(Q(Q,Q1)/100)
1810 B1=INT((Q(Q,Q1)-100*K1)*.1)
1820 IF B1=0 THEN 1850 ELSE GOSUB 2470:B2=INT(1+RND*8):GOSUB 2470:B3=INT(1+RND*8)
1830 IF A$(B2,B3)="*" THEN ST=ST-1
1840 IF A$(B2,B3)="E" THEN 1820 ELSE A$(B2,B3)="B"
1850 IF K1=0 THEN 1880 ELSE FOR A=1 TO K1:A$(P(A),N(A))=" ":GOSUB 2470:P(A)=INT(1+RND*8):GOSUB 2470:N(A)=INT(1+RND*8)
1860 IF A$(P(A),N(A))<>" " THEN 1850 ELSE A$(P(A),N(A))="K"
1870 IF O1=0 THEN GOSUB 2470:KE(A)=INT(150+RND*351):GOSUB 2470:KT(A)=INT(RND*3+.5):NEXT
1880 RETURN
1890 FOR A=Q-1 TO Q+1:FOR B7=Q1-1 TO Q1+1
1900 IF Q(A,B7)>=100 THEN C$="Yellow":RETURN
1910 NEXT:NEXT:RETURN
1920 REM *** REPAIR DAMAGE ***
1930 FOR A=0 TO 10:IF D(A)>0 AND D(A)<=TA THEN PRINT D$(A);" repair completed, ";N$:SP(A)=SP1(A):D(A)=0
1940 D(A)=D(A)-TA:IF D(A)<0 THEN D(A)=0
1950 IF LEN(STR$(D(A)))>12 THEN D(A)=0:PRINT D$(A);" repair completed, ";N$:SP(A)=SP1(A)
1960 NEXT:SD=SD+TA:MS=MS+(TA/2):TA=0:IF MS>=1 THEN GOSUB 2470:MS=RND*1
1970 REM *** KLINGONS  MOVE / FIRE ***
1980 IF K1=0 THEN 310 ELSE B77=1:PRINT:FOR A=1 TO K1:FOR B7=1 TO 3:KD(A,B7)=KD(A,B7)-.1:IF KD(A,B7)>0 THEN KD(A,B7)=0:NEXT:NEXT:A=1
1990 IF KD(A,1)>0 AND KD(A,2)>0 AND KD(A,3)>0 THEN 2090
1995 GOSUB 2470:IF RND*2<.5 AND KD(A,1)<=0 THEN GOSUB 2530:GOTO 2080
2000 GOSUB 2470:RR=1.25+RND*.0833333:SR=SH/RR
2010 GOSUB 2470:IF RND*1.5<.5 AND KT(A)>0 AND KD(A,3)<=0 THEN PRINT CHR$(7);"WARNING:  Klingon Torpedoes":KT(A)=KT(A)-1:GOSUB 2470:H=INT(200*(1+RND*.5)+.5):GOTO 2040
2020 IF KD(A,2)>0 THEN 1990 ELSE GOSUB 2470:X=INT(RND*KE(A)+KE(A)/10):GOSUB 2470:H=INT(4*(1+RND*.2-.4)*X/((P(A)-S)^2+(N(A)-S1)^2)^.34-X/10+.5)
2030 PRINT F$;"Klingons firing Phasers.":FOR B7=1 TO 3:PRINT CHR$(7);:NEXT
2040 IF D(10)=0 THEN H=H*INT(SH/30+.5):E=E-H:IF H<1 THEN PRINT " No damage sustained, ";N$ ELSE PRINT H;"units drained from Reserve, ";N$:X=H
2050 IF E<=20 THEN 2160 ELSE GOSUB 2470:B7=(H/E)*200:IF D(10)>0 THEN PRINT "Shields not functioning, ";M$;"!":GOTO 2100
2060 IF H>E/2 THEN PRINT "WARNING:  Shield overload!":DA=DA-1
2070 IF B7>=SR THEN 2100 ELSE PRINT " Shields still holding, ";N$
2080 IF B77=0 OR A>=K1 THEN B77=0:GOTO 310
2090 A=A+1:GOTO 1990
2100 GOSUB 2470:B7=INT(RND*13-1):IF B7>10 THEN B7=10 ELSE IF B7<0 THEN B7=0
2110 IF B7=10 THEN PRINT " ";E$;"Shields collapsing, ";M$;"!":GOTO 2130
2120 PRINT " ";E$;D$(B7);" damaged by the hit.":IF D(0)>0 THEN MS=0
2130 GOSUB 2470:D(B7)=D(B7)+RND*4:PRINT " Time to repair -";STR$(D(B7));".":SP(B7)=SP1(B7):DA=DA+1:IF DA<1 THEN 2100
2140 GOSUB 2470:IF RND*(X/H*H/100)>=1+D(10)/10 AND DA<5 THEN 2100 ELSE IF AK=10 THEN 310 ELSE DA=0:GOTO 2080
2150 REM *** ENTERPRISE DESTROYED ***
2160 PRINT:PRINT:PRINT CHR$(7);"The Enterprise has been destroyed.":LI=1:IF K<1 THEN 2240
2170 PRINT:PRINT "There are still";K;"Klingon Battle Cruisers."
2180 PRINT:PRINT "Better luck next time, ";N$
2190 R=((K2-K)/K2)*100:R=R-(SD1-SD):R=R-10*B8:R=R-K:R=R+(K2-K):R=R*(K2/50)
2200 PRINT:PRINT "Your rating was";STR$(R);"%, ";N$
2210 IF LI=0 THEN PRINT:PRINT:PRINT "You have been offered another commission, ";N$:INPUT "Will you accept it";P$
2220 P$=LEFT$(P$,1):IF P$="y" OR P$="Y" THEN RUN ELSE PRINT:PRINT:END
2230 REM *** MISSION ACCOMPLISHED ***
2240 PRINT:PRINT:PRINT "Congratulations, ";M$;"!!!  You have destroyed"
2250 PRINT:PRINT "all";K2;"Klingon Battle Cruisers!!!":GOTO 2190
2260 PRINT:PRINT F$;"We are on a Collision Course with a Star.":PRINT
2270 GOSUB 2470:FOR A=150-(10*I2+RND*(10*I2)) TO 4000+(30*I2+RND*(30*I2)) STEP I1*50+RND*(I1*200)
2280 PRINT F$;"Outer Hull Temperature:";A;"degrees and rising.":FOR B7=1 TO 700:NEXT B7
2290 NEXT:GOSUB 2470:PRINT:PRINT CHR$(7);"FIZZLE!  The Enterprise burned up";5000+RND*10000;"miles from the Star.":GOTO 2160
2300 REM *** COLLISION COURSE  W / KLINGON ***
2310 PRINT:PRINT F$;"We appear to be on a Collision Course with a Klingon Battle Cruiser."
2320 GOSUB 2470:FOR A=2500*I2+RND*(5000*I2) TO 250*I2+RND*(500*I2) STEP -(2500+RND*(I1*100))
2330 PRINT F$;"Distance to Klingon Cruiser:";A;"miles and closing.":FOR B7=1 TO 700:NEXT B7:NEXT A
2340 IF RND*8+I2/2>I1 THEN PRINT:PRINT F$;"The Klingons have outmaneuvered us, ";N$:O1=1:GOSUB 1850:GOTO 1930
2350 FOR AR=1 TO K1:IF P(AR)=S AND N(AR)=S1 THEN AK=AR
2360 NEXT
2370 K=K-1:K1=K1-1:PRINT CHR$(7):A$(P(AK),N(AK))=" ":FOR AF=AK TO K1:P(AF)=P(AF+1):N(AF)=N(AF+1):T(AF)=T(AF+1):T1(AF)=T1(AF+1):KE(AF)=KE(AF+1):NEXT:T3=T3-1:T2=T2-1:Q(Q,Q1)=Q(Q,Q1)-100:K1=K1-1
2380 IF SH>=50 AND E>2000 AND D(10)=0 THEN GOSUB 2470:SH=INT(SH*.5):E=E-2000:A=10:DA=-2:B7=10:AK=10:GOTO 2110
2390 PRINT "You know, that was awfully stupid of you --SIR.":K=K-1:GOTO 2160
2400 REM *** COLLISION COURSE  W / STARBASE ***
2410 PRINT:PRINT F$;"We are going to collide with our Starbase!":PRINT
2420 GOSUB 2470:FOR A=2500*I2+RND*(5000*I2) TO 250*I2+RND*(500*I2) STEP -(2500+RND*(I1*100))
2430 PRINT F$;"Distance to Starbase:";A;"miles and closing.":FOR B7=1 TO 700:NEXT B7:NEXT A:PRINT CHR$(7):B1=0:B=B-1:B8=B8+1
2440 PRINT F$;"We have just destroyed our own Starbase, ";N$:Q(Q,Q1)=Q(Q,Q1)-10
2450 IF SH>=25 AND E>1000 AND D(10)=0 THEN GOSUB 2470:SH=INT(SH*.75):E=E-1000:A=10:DA=-1:B7=10:GOTO 2110
2460 PRINT "BOOM!   That really wasn't too bright, ";M$;"!":GOTO 2160
2470 RANDOMIZE(VAL(RIGHT$(TIME$,2))):RETURN
2480 PRINT RIGHT$(STR$(Q2),1);",";RIGHT$(STR$(Q3),1);:RETURN
2490 GOSUB 2470:MS=MS-RND*.3:GOSUB 2470:IF RND*1.1>W3/10 THEN RETURN ELSE PRINT E$;"Engines superheating, ";M$;"!":DA=12:A=0:GOSUB 2470:A1=INT(5+RND*5):IF A1>=W3 THEN Q=Q2:Q1=Q3
2500 SS=0:FOR A=1 TO A1:BEEP:IF A=7 THEN PRINT E$;"The engines canna take much more of this, ";M$;"!"
2510 FOR B7=1 TO 500:NEXT B7:NEXT A:B7=0:GOSUB 2470:IF RND*1.2>W3/10 THEN 2520 ELSE K=K-K1:B8=B8+B1:GOTO 1720
2520 GOSUB 2470:D(0)=D(0)+RND*3:PRINT:PRINT E$;"I'm shuttin' down the Warp drive, "N$:PRINT " It'll take at least";INT(D(0)*100)/100;"Stardates to repair!"
2530 A$(P(A),N(A))=" ":GOSUB 2470:P(A)=INT(1+RND*8):GOSUB 2470:N(A)=INT(1+RND*8)
2540 IF A$(P(A),N(A))<>" " THEN 2530 ELSE A$(P(A),N(A))="K":RETURN
2550 IF Q(A,B7)=0 THEN G(A,B7)=.1
2560 RETURN
2570 IF Q(Q,Q1)=0 THEN G(Q,Q1)=.1:Q(Q,Q1)=.1
2580 IF K1>0 THEN RETURN ELSE FOR A=1 TO 8:FOR B7=1 TO 8
2590 IF A$(A,B7)="K" THEN A$(A,B7)=" "
2600 NEXT:NEXT:IF B1=1 THEN GOSUB 2470:PRINT "This is Starbase";INT(1+RND*B);"at Sector ";:Q2=B2:Q3=B3:GOSUB 2480:PRINT ":  Good job, Enterprise!"
2610 RETURN

