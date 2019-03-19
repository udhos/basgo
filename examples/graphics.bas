100 screen 1000
200 line (1,1)-(640,480)
210 for i=240 to 480: line (400,240)-(640,i): next
220 line (10,100)-(300,300),2,bf,2
900 _goimport("time")
910 _godecl("func sleep() { time.Sleep(3*time.Second) }")
920 _goproc("sleep")
930 end
