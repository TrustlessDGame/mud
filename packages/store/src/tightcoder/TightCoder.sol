// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

import { SchemaType } from "@latticexyz/schema-type/src/solidity/SchemaType.sol";
import { Slice, SliceLib } from "../Slice.sol";

library TightCoder {
  /**
   * @dev Copies the array to a new bytes array,
   * tightly packing it using the given size per element (in bytes)
   */
  function encode(
    bytes32[] memory array,
    uint256 elementSize,
    uint256 shiftLeftBits
  ) internal pure returns (bytes memory data) {
    uint256 arrayLength = array.length;
    uint256 packedLength = array.length * elementSize;
    data = new bytes(packedLength);

    /// @solidity memory-safe-assembly
    assembly {
      for {
        let i := 0
        // Skip array length
        let fromPointer := add(array, 0x20)
        let toPointer := add(data, 0x20)
      } lt(i, arrayLength) {
        // Loop until we reach the end of the array
        i := add(i, 1)
        // Increment array pointer by one word
        fromPointer := add(fromPointer, 0x20)
        // Increment packed pointer by one element size
        toPointer := add(toPointer, elementSize)
      } {
        mstore(toPointer, shl(shiftLeftBits, mload(fromPointer))) // pack one array element
      }
    }
  }

  /**
   * @dev Unpacks the slice to a new memory location
   * and lays it out like a memory array with the given size per element (in bytes)
   * @return array a generic array, needs to be casted to the expected type using assembly
   */
  function decode(
    Slice packedSlice,
    uint256 elementSize,
    bool leftAligned
  ) internal pure returns (bytes32[] memory array) {
    uint256 packedPointer = packedSlice.pointer();
    uint256 packedLength = packedSlice.length();
    uint256 padLeft = leftAligned ? 0 : 256 - elementSize * 8;
    // Array length (number of elements)
    uint256 arrayLength;
    unchecked {
      arrayLength = packedLength / elementSize;
    }

    // TODO temporary check to catch bugs, either remove it or use a custom error
    // (see https://github.com/latticexyz/mud/issues/444)
    if (packedLength % elementSize != 0) {
      revert("unpackToArray: packedLength must be a multiple of elementSize");
    }

    /// @solidity memory-safe-assembly
    assembly {
      // Allocate a word for each element, and a word for the array's length
      let allocateBytes := add(mul(arrayLength, 32), 0x20)
      // Allocate memory and update the free memory pointer
      array := mload(0x40)
      mstore(0x40, add(array, allocateBytes))

      // Store array length
      mstore(array, arrayLength)

      for {
        let i := 0
        let arrayCursor := add(array, 0x20) // skip array length
        let packedCursor := packedPointer
      } lt(i, arrayLength) {
        // Loop until we reach the end of the array
        i := add(i, 1)
        arrayCursor := add(arrayCursor, 0x20) // increment array pointer by one word
        packedCursor := add(packedCursor, elementSize) // increment packed pointer by one element size
      } {
        mstore(arrayCursor, shr(padLeft, mload(packedCursor))) // unpack one array element
      }
    }
  }
}
