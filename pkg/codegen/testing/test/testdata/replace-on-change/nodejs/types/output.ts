// *** WARNING: this file was generated by test. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";
import * as utilities from "../utilities";

import {Cat, Dog} from "..";

/**
 * A toy for a dog
 */
export interface Chew {
    owner?: Dog;
}

/**
 * A Toy for a cat
 */
export interface Laser {
    animal?: Cat;
    batteries?: boolean;
    light?: number;
}

export interface Rec {
    rec1?: outputs.Rec;
}

/**
 * This is a toy
 */
export interface Toy {
    associated?: outputs.Toy;
    color?: string;
    wear?: number;
}

