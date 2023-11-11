: fizz?  3 mod 0 = dup if ." Fizz " then ;
: buzz?  5 mod 0 = dup if ." Buzz " then ;
: another? dup fizz? swap buzz? or invert ;
: something 25 1 do cr i izz-uzz? if i . then loop ;
something
