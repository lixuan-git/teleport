/* eslint-disable */
// @generated by protobuf-ts 2.9.3 with parameter eslint_disable,add_pb_suffix,server_grpc1,ts_nocheck
// @generated from protobuf file "teleport/userpreferences/v1/access_graph.proto" (package "teleport.userpreferences.v1", syntax proto3)
// tslint:disable
// @ts-nocheck
//
// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
/**
 * AccessGraphUserPreferences is the user preferences for Access Graph.
 *
 * @generated from protobuf message teleport.userpreferences.v1.AccessGraphUserPreferences
 */
export interface AccessGraphUserPreferences {
    /**
     * has_been_redirected is true if the user has already been redirected to the Access Graph
     * on login, after having signed up for a trial from the Teleport Policy page.
     *
     * @generated from protobuf field: bool has_been_redirected = 1;
     */
    hasBeenRedirected: boolean;
}
// @generated message type with reflection information, may provide speed optimized methods
class AccessGraphUserPreferences$Type extends MessageType<AccessGraphUserPreferences> {
    constructor() {
        super("teleport.userpreferences.v1.AccessGraphUserPreferences", [
            { no: 1, name: "has_been_redirected", kind: "scalar", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<AccessGraphUserPreferences>): AccessGraphUserPreferences {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.hasBeenRedirected = false;
        if (value !== undefined)
            reflectionMergePartial<AccessGraphUserPreferences>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: AccessGraphUserPreferences): AccessGraphUserPreferences {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* bool has_been_redirected */ 1:
                    message.hasBeenRedirected = reader.bool();
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: AccessGraphUserPreferences, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* bool has_been_redirected = 1; */
        if (message.hasBeenRedirected !== false)
            writer.tag(1, WireType.Varint).bool(message.hasBeenRedirected);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message teleport.userpreferences.v1.AccessGraphUserPreferences
 */
export const AccessGraphUserPreferences = new AccessGraphUserPreferences$Type();
