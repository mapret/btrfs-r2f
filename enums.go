package btrfsr2f

type sendCommand uint16

const (
	BTRFS_SEND_C_UNSPEC       sendCommand = iota // 0
	BTRFS_SEND_C_SUBVOL                          // 1
	BTRFS_SEND_C_SNAPSHOT                        // 2
	BTRFS_SEND_C_MKFILE                          // 3
	BTRFS_SEND_C_MKDIR                           // 4
	BTRFS_SEND_C_MKNOD                           // 5
	BTRFS_SEND_C_MKFIFO                          // 6
	BTRFS_SEND_C_MKSOCK                          // 7
	BTRFS_SEND_C_SYMLINK                         // 8
	BTRFS_SEND_C_RENAME                          // 9
	BTRFS_SEND_C_LINK                            // 10
	BTRFS_SEND_C_UNLINK                          // 11
	BTRFS_SEND_C_RMDIR                           // 12
	BTRFS_SEND_C_SET_XATTR                       // 13
	BTRFS_SEND_C_REMOVE_XATTR                    // 14
	BTRFS_SEND_C_WRITE                           // 15
	BTRFS_SEND_C_CLONE                           // 16
	BTRFS_SEND_C_TRUNCATE                        // 17
	BTRFS_SEND_C_CHMOD                           // 18
	BTRFS_SEND_C_CHOWN                           // 19
	BTRFS_SEND_C_UTIMES                          // 20
	BTRFS_SEND_C_END                             // 21
)

/*type tlvType uint16

const (
	BTRFS_TLV_U8       tlvType = iota // 0
	BTRFS_TLV_U16                     // 1
	BTRFS_TLV_U32                     // 2
	BTRFS_TLV_U64                     // 3
	BTRFS_TLV_BINARY                  // 4
	BTRFS_TLV_STRING                  // 5
	BTRFS_TLV_UUID                    // 6
	BTRFS_TLV_TIMESPEC                // 7
)*/

type sendAttribute uint16

const (
	BTRFS_SEND_A_UNSPEC         sendAttribute = iota // 0
	BTRFS_SEND_A_UUID                                // 1
	BTRFS_SEND_A_CTRANSID                            // 2
	BTRFS_SEND_A_INO                                 // 3
	BTRFS_SEND_A_SIZE                                // 4
	BTRFS_SEND_A_MODE                                // 5
	BTRFS_SEND_A_UID                                 // 6
	BTRFS_SEND_A_GID                                 // 7
	BTRFS_SEND_A_RDEV                                // 8
	BTRFS_SEND_A_CTIME                               // 9
	BTRFS_SEND_A_MTIME                               // 10
	BTRFS_SEND_A_ATIME                               // 11
	BTRFS_SEND_A_OTIME                               // 12
	BTRFS_SEND_A_XATTR_NAME                          // 13
	BTRFS_SEND_A_XATTR_DATA                          // 14
	BTRFS_SEND_A_PATH                                // 15
	BTRFS_SEND_A_PATH_TO                             // 16
	BTRFS_SEND_A_PATH_LINK                           // 17
	BTRFS_SEND_A_FILE_OFFSET                         // 18
	BTRFS_SEND_A_DATA                                // 19
	BTRFS_SEND_A_CLONE_UUID                          // 20
	BTRFS_SEND_A_CLONE_CTRANSID                      // 21
	BTRFS_SEND_A_CLONE_PATH                          // 22
	BTRFS_SEND_A_CLONE_OFFSET                        // 23
	BTRFS_SEND_A_CLONE_LEN                           // 24
)
