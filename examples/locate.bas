100 screen 0
110 brick$="x"
200 for i=1 to 11
300 for c=i to 81-i
310 if c<80 then locate 26-i,c: print brick$;
320 next
400 for r=26-i to i step -1
410 if r<25 then locate r,81-i: print brick$;
420 next
500 for c=81-i to i step -1
510 locate i,c: print brick$;
520 next
600 for r=i to 26-i
610 locate r,i: print brick$;
620 next
900 next
990 print input$(1)
