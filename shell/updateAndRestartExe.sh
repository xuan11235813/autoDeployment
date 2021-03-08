exeDirectory="radarSystemRaw"
runFileTemp="outMultipleTemp.lexe"
runFile="outMultiple.lexe"
# kill all screens
killall screen
# rename the executable files
mv $runFileTemp $exeDirectory/$runFile
chmod +x $exeDirectory/$runFile
# restart the program
screen -S targetSession -d -m bash
screen -r targetSession -X stuff "cd $exeDirectory"$(echo -ne '\015')
screen -r targetSession -X stuff "sudo ./$runFile"$(echo -ne '\015')

sleep 1

# test using pgrep
if ! pgrep "out" > /dev/null
then
    echo "no program started"
else
    echo "program starts again"
fi