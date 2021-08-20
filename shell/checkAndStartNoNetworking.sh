configFileName="vec.config"
runFile="outMultiple.lexe"
runDirectory="radarSystem"
# prog=$(top -n2 | grep -i out | awk '{print $2}')
# if test -z "$prog"
if ! pgrep "out" > /dev/null
then
    echo "start the normal program"
    killall screen
	screen -S targetSession -d -m bash
	screen -r targetSession -X stuff "cd $runDirectory"$(echo -ne '\015')
	screen -r targetSession -X stuff "sudo ./$runFile 5"$(echo -ne '\015')
else
    echo "program is running"
fi