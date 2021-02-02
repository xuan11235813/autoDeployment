configFileName="vec.config"
runFile1="outMultiple.lexe"
runFile2="outMultipleTest.lexe"
runDirectory="radarSystemRaw"
# prog=$(top -n2 | grep -i out | awk '{print $2}')
# if test -z "$prog"
if ! pgrep "out" > /dev/null
then
    if test -f $runDirectory/$configFileName
    then
        echo "start the normal program"
        killall screen
		screen -S targetSession -d -m bash
		screen -r targetSession -X stuff "cd $runDirectory"$(echo -ne '\015')
		screen -r targetSession -X stuff "sudo ./$runFile1"$(echo -ne '\015')
    else
        echo "start the test program"
        killall screen
		screen -S targetSession -d -m bash
		screen -r targetSession -X stuff "cd $runDirectory"$(echo -ne '\015')
		screen -r targetSession -X stuff "sudo ./$runFile2"$(echo -ne '\015')
    fi
else
    echo "program is running"
fi