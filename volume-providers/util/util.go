/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018, 2019 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpcvolume

import (
	"fmt"
	"log"
	"time"
)

// Get execution time of a function
func TimeTracker(functionName string, start time.Time) {
	elapsed := time.Since(start)

	log.Println(fmt.Sprintf("Time taken by function %s is %s", functionName, elapsed))
}
