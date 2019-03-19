100 screen 1000
200 line (1,1)-(640,480)
210 _goproc("sleep1")
220 cls
230 for i=240 to 480: line (400,240)-(640,i): next
240 line (10,100)-(300,300),2,b,2
250 line (20,110)-(290,290),2,bf,2
900 _goimport("time")
910 _godecl("func sleep1() { time.Sleep(1*time.Second) }")
920 _godecl("func sleep3() { time.Sleep(3*time.Second) }")
930 _goproc("sleep3")
940 end
