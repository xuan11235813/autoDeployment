runFile2="outMultipleTest.lexe"
runDirectory="radarSystemRaw"

if ! pgrep "out" > /dev/null
then
    if ! test -f $runDirectory/$runFile2
    then
        echo "start the test program"
        killall screen
		screen -S targetSession -d -m bash
		screen -r targetSession -X stuff "cd $runDirectory"$(echo -ne '\015')
		screen -r targetSession -X stuff "sudo ./$runFile2"$(echo -ne '\015')
    fi
else
    echo "program is running"
fi