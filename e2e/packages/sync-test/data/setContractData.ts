import { Page } from "@playwright/test";
import { Data } from "./types";
import { encodeTestData } from "./encodeTestData";
import { callWorld } from "./callWorld";
import { stringToBytes16 } from "@latticexyz/utils";
import { toHex } from "viem";

/**
 * Writes contract data by calling `world.setRecord` via the client
 */
export async function setContractData(page: Page, data: Data) {
  const encodedData = encodeTestData(data);
  const promises: Promise<unknown>[] = [];
  for (const [table, records] of Object.entries(encodedData)) {
    for (const record of records) {
      const promise = callWorld(page, "setRecord(bytes16,bytes16,bytes32[],bytes)", [
        // TODO: add support for multiple namespaces after https://github.com/latticexyz/mud/issues/994 is resolved
        toHex(stringToBytes16("")),
        toHex(stringToBytes16(table)),
        record.key,
        record.value,
      ]);

      // Wait for transactions to be confirmed
      promises.push(promise);
    }
  }
  return Promise.all(promises);
}
