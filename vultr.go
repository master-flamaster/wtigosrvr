Last login: Wed Jul 16 07:48:52 on ttys000
dima@MacBook-Pro-dima ~ % # Процессы, отсортированные по CPU
top -o %CPU
# Или более дружелюбный вариант
htop          # sudo apt install htop

# Какой процесс занял порт 8080
sudo lsof -i :8080
# Все открытые TCP/UDP-сокеты
sudo ss -tulpn
zsh: command not found: #
invalid argument -o: %CPU
top usage: top
		[-a | -d | -e | -c <mode>]
		[-F | -f]
		[-h]
		[-i <interval>]
		[-l <samples>]
		[-ncols <columns>]
		[-o <key>] [-O <secondaryKey>]
			keys: pid (default), command, cpu, cpu_me, cpu_others, csw,
				time, threads, ports, mregion, mem, rprvt, purg, vsize, vprvt,
				kprvt, kshrd, pgrp, ppid, state, uid, wq, faults, cow, user,
				msgsent, msgrecv, sysbsd, sysmach, pageins, boosts, instrs, cycles
		[-R | -r]
		[-S]
		[-s <delay>]
		[-n <nprocs>]
		[-stats <key(s)>]
		[-pid <processid>]
		[-user <username>]
		[-U <username>]
		[-u]

zsh: command not found: #
zsh: command not found: htop
zsh: command not found: #
Password:
Sorry, try again.
Password:
zsh: command not found: #
sudo: ss: command not found
dima@MacBook-Pro-dima ~ % htop          # sudo apt install htop
zsh: command not found: htop
dima@MacBook-Pro-dima ~ % sudo lsof -i :8080
dima@MacBook-Pro-dima ~ % sudo ss -tulpn
sudo: ss: command not found
dima@MacBook-Pro-dima ~ % # Процессы, отсортированные по CPU
top -o %CPU
# Или более дружелюбный вариаhtop          # sudo apt install htop

# Какой процесс занял порт 8080
sudo lsof -i :8080
# Все открытые TCP/UDP-сокеты
sudo ss -tulpn
