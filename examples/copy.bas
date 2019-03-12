100 from$=environ$("FROM")
110 to$=environ$("TO")
120 print "will copy bytes FROM=";from$;" TO=";to$
130 if len(from$)<1 then print "please set the env var FROM": end
140 if len(to$)<1 then print "please set the env var TO": end
150 print "opening files"
160 open from$ for input as 1
170 open to$ for output as 2
180 print "copying bytes FROM=";from$;" TO=";to$
190 while not eof(1)
200 c$=input$(1,#1)
210 print #2,c$;
220 count%=count%+len(c$)
230 wend 
240 print count%;"bytes copied"
