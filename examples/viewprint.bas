
100 screen 0
110 locate 1,1: print "---screen first line";
120 locate 25,1: print "---screen last line";
130 vp1=1+5:vp2=25-5
140 locate vp1-1,1: print "###before view print";
150 locate vp2+1,1: print "###after view print";
160 view print vp1 to vp2
170 locate 1,1: print "hit ENTER to clear screen"
180 i$=input$(1)
190 cls
200 print "screen cleared"
210 print "hit a key or q to quit"
220 i$=input$(1)
230 while i$ <> "q"
240 print "key: "; i$
250 print "line:";csrlin;" - hit another key, c to clear screen, or q to quit - ";
260 i$=input$(1)
270 if i$="c" then cls:print "screen cleared"
280 wend
