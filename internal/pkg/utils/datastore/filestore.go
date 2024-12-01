package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

// store data in file

// expected data in yaml:
// cat <<EOF | kubectl apply -f -
// apiVersion: gateway.envoyproxy.io/v1alpha1
// kind: BackendTrafficPolicy
// metadata:
//   name: policy-httproute
// spec:
//   targetRefs:
//   - group: gateway.networking.k8s.io
//     kind: HTTPRoute
//     name: http-ratelimit
//   rateLimit:
//     type: Global
//     global:
//       rules:
//       - clientSelectors:
//         - sourceCIDR:
//           value: 0.0.0.0/0
//           type: Distinct
//         - headers:
//           - name: x-user-id
//             value: one
//           - type: Distinct
//             name: x-user-id
//           - name: x-user-id
//             value: admin
//             invert: true

//         limit:
//           requests: 3
//           unit: Hour
// EOF

// rateLimitHttpRoute: http-ratelimit
// rateLimitType: Global
// rateLimitRules[]

func readJSONFile(filepath string) ([]byte, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    byteValue, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }
    return byteValue, nil
}

func convertJSONToYAML(jsonData []byte) ([]byte, error) {
    // convert json to yaml
    yamlData, err := yaml.JSONToYAML(jsonData)
    if err != nil {
        return nil, err
    }
    return yamlData, nil
}

// func generateBackendTrafficPolicy(yamlData []byte) (*egv1a1.BackendTrafficPolicy, error) {
//     btpResource := &egv1a1.BackendTrafficPolicy{}
//     yamlInput := string(yamlData)

//     // get name from yaml with key rateLimitHttpRoute
//     // get rateLimitType from yaml with key rateLimitType
//     name :=  yamlInput["rateLimitHttpRoute"]
//     rateLimitType := yamlInput["rateLimitType"]
//     rateLimitRules := yamlInput["rateLimitRules"]

//     // set values in btpResource
//     // targetRef := btpResource.Spec.TargetRefs[0]
//     // targetRef.Group = "gateway.networking.k8s.io"
//     // targetRef.Kind = "HTTPRoute"
//     // targetRef.Name = name

//     // btpTargetRef := gwapiv1a2.LocalPolicyTargetReferenceWithSectionName{
//     //                     LocalPolicyTargetReference: gwapiv1a2.LocalPolicyTargetReference{
//     //                         Group: "gateway.networking.k8s.io",
//     //                         Kind: "HTTPRoute",
//     //                         Name: 
//     //                     },
//     // }


//     // btpResource.Spec.TargetRefs = append(btpResource.Spec.TargetRefs, targetRef)
//     // btpResource.Spec.RateLimit.Type = egv1a1.RateLimitType(rateLimitType)
//     // btpResource.Spec.RateLimit.Global.Rules = rateLimitRules

// }

type BTPData struct {
    RateLimitHttpRoute string `json:"rateLimitHttpRoute"`
    RateLimitType egv1a1.RateLimitType `json:"rateLimitType"`
    RateLimitRules []egv1a1.RateLimitRule `json:"rateLimitRules"`
}

func main() {
    // get path of the file
    cmd, err := os.Getwd()
    if err != nil {
        log.Println("Error getting current working directory: ", err)
    }
    fileName := "btp.json"
    fileDir := ""
    filePath := filepath.Join(cmd, fileDir, fileName)
    fmt.Println("File path: ", filePath)


    // read json data from file    
    // jsonData, err := readJSONFile(filePath)
    // if err != nil {
    //     log.Println("Error reading JSON file: ", err)
    // }

    // yamlData, err := convertJSONToYAML(jsonData)
    // if err != nil {
    //     log.Println("Error converting JSON to YAML: ", err)
    // }

    // //  print yaml data
    // // fmt.Println(string(yamlData))
    // // read section 'rateLimitType' from yaml
    // var yamlMap map[string]interface{}
    // // err = yaml.


    // generate BackendTrafficPolicy object
    // btpResource, _ := generateBackendTrafficPolicy(yamlData)
    // fmt.Println("bptResource: ", *btpResource)

    jsonDataBytes, err := os.ReadFile(filePath)
    if err != nil {
        log.Println("Error reading JSON file: ", err)
    }
    // var jsonData []interface{}
    // err = json.Unmarshal(jsonDataBytes, &jsonData)
    // if err != nil {
    //     log.Println("Error unmarshalling JSON data: ", err)
    // }

    // fmt.Println("jsonData: ", jsonData.rateLimitHttpRoute, jsonData.rateLimitType)
    // fmt.Println("jsonData: ", jsonData[0].(map[string]interface{})["rateLimitRules"])

    // var rateLimitType egv1a1.RateLimitType
    // rateLimitType := egv1a1.RateLimitType(jsonData[0].()

    var btpData []BTPData
    err = json.Unmarshal(jsonDataBytes, &btpData)
    if err != nil {
        log.Println("Error unmarshalling JSON data: ", err)
    }
    fmt.Println("btpData: ", btpData[0].RateLimitHttpRoute, btpData[0].RateLimitType, btpData[0].RateLimitRules)


    // Convert Go struct to yaml
    yamlData, err := yaml.Marshal(btpData)
    if err != nil {
        log.Println("Error converting Go struct to YAML: ", err)
    }

    fmt.Println(string(yamlData))
}