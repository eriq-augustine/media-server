package util;

import (
   "crypto/sha1"
   "encoding/hex"
);

// Get the hex SHA1 string.
func SHA1Hex(val string) string {
   hash := sha1.New();
   hash.Write([]byte(val));
   return hex.EncodeToString(hash.Sum(nil));
}
