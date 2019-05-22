// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

#include <stdio.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>
#include "unity.h"
#include "../../uplink-cgo.h"

UplinkRef NewTestUplink(char **);

void TestCreateBucket(void)
{
    char *_err = "";
    char **err = &_err;
    char *satelliteAddr = getenv("SATELLITEADDR");
    gvAPIKey apiKey = ParseAPIKey(getenv("APIKEY"), err);

    pbProjectOptions projectOpts = STORJ__LIBUPLINK__PROJECT_OPTIONS__INIT;
    // NB: empty encryption key
    uint8_t encryptionKey[32];
    memset(encryptionKey, '\0', sizeof(encryptionKey));
    projectOpts.encryption_key.data = encryptionKey;
    projectOpts.encryption_key.len = sizeof(encryptionKey);

    gvProjectOptions *optsValue = malloc(sizeof(gvProjectOptions));
    optsValue->Type = ProjectOptionsType;
    protoToGoValue(&projectOpts, optsValue, err);
    TEST_ASSERT_EQUAL_STRING("", *err);

    UplinkRef uplinkRef = NewTestUplink(err);
    TEST_ASSERT_EQUAL_STRING("", *err);
    TEST_ASSERT_NOT_EQUAL(0, uplinkRef);

    ProjectRef projectRef = OpenProject(uplinkRef, satelliteAddr, apiKey.Ptr, *optsValue, err);
    TEST_ASSERT_EQUAL_STRING("", *err);

    pbBucketConfig bucketCfg = STORJ__LIBUPLINK__BUCKET_CONFIG__INIT;
    bucketCfg.path_cipher = 0;
//    bucketCfg.encryption_parameters = ;
//    bucketCfg.redundancy_scheme = ;
    bucketCfg.segment_size = 1024;

    // NB: dev defaults (maybe factor out into a lib helper)
    pbRedundancyScheme scheme = STORJ__LIBUPLINK__REDUNDANCY_SCHEME__INIT;
    scheme.algorithm = 1;
    scheme.share_size = 1024;
    scheme.required_shares = 4;
    scheme.optimal_shares = 8;
    scheme.total_shares = 10;

//    TEST_ASSERT_EQUAL_STRING("", *err);
//    TEST_ASSERT_EQUAL_STRING("testbucket", bucket.Name);
}
