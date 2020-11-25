# This file is used to create the btrfs-stream files that are used for integration testing
#
# Preconditions:
#  - btrfs file system mounted on /mnt/data
#  - btrfs subvolume /mnt/data/main
#  - btrfs subvolume /mnt/data/snapshots
#  - Linux environment (obviously)
#

TEXT1="The quick brown fox jumps over the lazy dog."
TEXT2="abcdefghijklmnopqrstuvwxyz\nABCDEFGHIJKLMNOPQRSTUVWXYZ\n0123456789"
PWD2=$(pwd)
USER=$(whoami)

cd /mnt/data/main
mkdir dir1
mkdir dir1/dir2
echo $TEXT1 > qbf.txt
echo $TEXT2 > letters.txt
echo $TEXT1 > dir1/qbf.txt
echo $TEXT2 > dir1/dir2/letters.txt
ln -s dir1/qbf.txt qbf_symlink.txt
ln dir1/dir2/letters.txt letters_hardlink.txt
sudo btrfs subvolume snapshot -r /mnt/data/main /mnt/data/snapshots/01
echo $PWD2/stream01.bin
sudo btrfs send /mnt/data/snapshots/01 -f $PWD2/stream01.bin
sudo chown $USER:$USER $PWD2/stream01.bin
find . -type f,l -exec md5sum {} \; | sort > $PWD2/stream01_files.txt

rm dir1/dir2/letters.txt
rmdir dir1/dir2
truncate -s10 dir1/qbf.txt
echo -n "___" | dd conv=notrunc of=letters.txt bs=1 seek=5
mv qbf.txt dir1/qbf2.txt
sudo btrfs subvolume snapshot -r /mnt/data/main /mnt/data/snapshots/02
sudo btrfs send -p /mnt/data/snapshots/01 /mnt/data/snapshots/02 -f $PWD2/stream02.bin
sudo chown $USER:$USER $PWD2/stream02.bin
find . -type f,l -exec md5sum {} \; | sort > $PWD2/stream02_files.txt

sudo btrfs subvolume delete /mnt/data/snapshots/01
sudo btrfs subvolume delete /mnt/data/snapshots/02
