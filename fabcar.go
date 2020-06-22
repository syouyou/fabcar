/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

 package main

 import (
	 "encoding/json"
	 "fmt"
	 "strconv"
	 "time"
	 "github.com/hyperledger/fabric-contract-api-go/contractapi"
 )
 
 // SmartContract provides functions for managing a HmContract
 type SmartContract struct {
	 contractapi.Contract
 }
 
 // 合同信息
 type HmContract struct {
	Id   string `json:"id"`
	ContractNo   string `json:"contractNo"`
	Amount string `json:"amount"`
	LoanDate string `json:"loanDate"`
	BusinessType string `json:"businessType"`
	Lender  Lender `json:"lender"`
	Borrower  Borrower `json:"borrower"`
 }
 
 type Lender struct {
	 Name   string `json:"name"`
	 LegalRepresentative  string `json:"legalRepresentative"`
	 Idcard   string `json:"idcard"`
	 Phone   string `json:"phone"`
	 Address   string `json:"address"`
 }
 
 type Borrower struct {
	 Name   string `json:"name"`
	 Gender  string `json:"gender"`
	 Idcard   string `json:"idcard"`
	 Phone   string `json:"phone"`
	 Address   string `json:"address"`
	 Age   string `json:"age"`
	 Nation   string `json:"nation"`
 }
 
 //业务证据信息
 type HmEvidence struct {
	 Content   string `json:"content"`
	 DataType  string `json:"dataType"`
	 FileName string `json:"fileName"`
	 EvidenceType  string `json:"evidenceType"`
	 BizId  string `json:"bizId"`
	 AddInfo  string `json:"addInfo"`
	//  class string `metadata:"class"`
	//  key string `metadata:"key"`
 }
 
 // QueryResult structure used for handling result of query
 type QueryResult struct {
	Key    string `json:"key"`
	Value string `json:"value"`
}

type QueryResultByPage struct {
	Size    string `json:"size"`
	BookMark string `json:"bookMark"`
	QueryResults []QueryResult `json:"queryResults"`
}

 type QueryHistoryResult struct {
	 TxId string `json:"tx_id"`
	 Value string `json:"value"`
	 IsDel string `json:"is_del"`
	 OnChainTime string `json:"on_chain_time"`
 }
 
 // InitLedger adds a base set of hmContracts to the ledger
 func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	 hmContracts := []HmContract{
		HmContract{ContractNo: "test",},
	 }
 
	 for i, hmContract := range hmContracts {
		 hmContractAsBytes, _ := json.Marshal(hmContract)
		 err := ctx.GetStub().PutState("hmContract"+strconv.Itoa(i), hmContractAsBytes)
 
		 if err != nil {
			 return fmt.Errorf("Failed to put to world state. %s", err.Error())
		 }
	 }
 
	 return nil
 }
 
 //1业务证据信息上传
 func (s *SmartContract) CreateHmEvidence(ctx contractapi.TransactionContextInterface, contentKey string, content string, dataType string, fileName string, evidenceType string, bizId string, addInfo string) error {
	 hmEvidence := HmEvidence{
		 Content:   content,
		 DataType:  dataType,
		 FileName: fileName,
		 EvidenceType:  evidenceType,
		 BizId: bizId,
		 AddInfo:addInfo,
	 }
 
	 hmEvidenceAsBytes, _ := json.Marshal(hmEvidence)
 
	 return ctx.GetStub().PutState(contentKey, hmEvidenceAsBytes)
 }

//2信息查询
func (s *SmartContract) QueryInfo(ctx contractapi.TransactionContextInterface, param string) ([]QueryResult, error) {
	queryIterator, err := ctx.GetStub().GetQueryResult(param)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	defer queryIterator.Close()
	results := make([]QueryResult, 0)
	for queryIterator.HasNext() {
		item, _ := queryIterator.Next()
		res := QueryResult{}
		res.Value=string(item.Value)
		res.Key=item.Key
		results= append(results, res)
	}

	return results, nil
}

//4信息查询--分页
func (s *SmartContract) QueryInfoByPage(ctx contractapi.TransactionContextInterface, queryString string, size string, bookmark string) (*QueryResultByPage, error) {
	var pageSize int32 = -1
	value, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		return nil, err
	}
	pageSize = int32(value)

	resultsIterator, metadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	
	queryResultByPage := new(QueryResultByPage)
	results := []QueryResult{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		
		temp := QueryResult{}
		temp.Key = queryResponse.Key
		temp.Value = string(queryResponse.Value)
		results= append(results, temp)
	}

	queryResultByPage.Size=fmt.Sprintf("%v", metadata.FetchedRecordsCount)
	queryResultByPage.BookMark = metadata.Bookmark
	queryResultByPage.QueryResults = results

	return queryResultByPage, nil
}

//3合同信息上传
func (s *SmartContract) CreateHmContract(ctx contractapi.TransactionContextInterface, key string, id string, contractNo string, amount string, loanDate string, businessType string, lender string, borrower string) error {
	lenderObj:=Lender{}
	if len(lender) != 0{
		err:=json.Unmarshal([]byte(lender),&lenderObj)
		if err!=nil{
			return fmt.Errorf("Failed to read from world state. %s", err.Error())
		}
	}	
	
	borrowerObj:=Borrower{}
	if len(borrower) != 0{
		err:=json.Unmarshal([]byte(borrower),&borrowerObj)
		if err!=nil{
			return fmt.Errorf("Failed to read from world state. %s", err.Error())
		}
	}
	
	hmContract := HmContract{
		ContractNo:   contractNo,
		Id:  id,
		Amount: amount,
		LoanDate:  loanDate,
		BusinessType:  businessType,
		Lender: lenderObj,
		Borrower: borrowerObj,
	}

	hmContractAsBytes, _ := json.Marshal(hmContract)

	return ctx.GetStub().PutState(key, hmContractAsBytes)
}

 // 根据key查询业务证据信息历史记录
 func (s *SmartContract) QueryHistory(ctx contractapi.TransactionContextInterface, content_key string) ([]QueryHistoryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(content_key)
		if err != nil {
				return nil, err
		}
		defer resultsIterator.Close()
		results := make([]QueryHistoryResult, 0)
		for resultsIterator.HasNext() {
				if queryResponse, err := resultsIterator.Next();err==nil{
						res := QueryHistoryResult{}
						res.TxId=queryResponse.TxId
						res.Value=string(queryResponse.Value)
						res.IsDel=strconv.FormatBool(queryResponse.IsDelete)
						res.OnChainTime=time.Unix(queryResponse.Timestamp.Seconds,0).Format("2006-01-02 15:04:05")
						results= append(results, res)
				}
				if err!=nil {
						return nil,err
				}
		}
		return results, nil
}
 
//  // 根据key查询业务证据信息
//  func (s *SmartContract) QueryHmEvidence(ctx contractapi.TransactionContextInterface, content_key string) (*HmEvidence, error) {
// 	hmBytes, err := ctx.GetStub().GetState(content_key)

// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
// 	}

// 	if hmBytes == nil {
// 		return nil, fmt.Errorf("%s does not exist", content_key)
// 	}

// 	hmEvidence := new(HmEvidence)
// 	_ = json.Unmarshal(hmBytes, hmEvidence)

// 	return hmEvidence, nil
// }
 
 
 
 //根据key查询合同信息
//  func (s *SmartContract) QueryHmContract(ctx contractapi.TransactionContextInterface, key string) (*HmContract, error) {
//  	hmContractAsBytes, err := ctx.GetStub().GetState(key)
 
//  	if err != nil {
//  		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
//  	}
 
//  	if hmContractAsBytes == nil {
//  		return nil, fmt.Errorf("%s does not exist", key)
//  	}
 
//  	hmContract := new(HmContract)
//  	_ = json.Unmarshal(hmContractAsBytes, hmContract)
 
//  	return hmContract, nil
//  }
 
//  func (s *SmartContract) QueryHmContractHis(ctx contractapi.TransactionContextInterface, caseNo string) ([]QueryHistoryResult, error) {
// 	 resultsIterator, err := ctx.GetStub().GetHistoryForKey(caseNo)
// 		 if err != nil {
// 				 return nil, err
// 		 }
// 		 defer resultsIterator.Close()
// 		 //results := []QueryResult{}
// 		 //results := make([]QueryResult, 0)
// 		 results := make([]QueryHistoryResult, 0)
// 		 for resultsIterator.HasNext() {
// 				 if queryResponse, err := resultsIterator.Next();err==nil{
// 						 res := QueryHistoryResult{}
// 						 res.TxId=queryResponse.TxId
// 						 res.Value=string(queryResponse.Value)
// 						 res.IsDel=strconv.FormatBool(queryResponse.IsDelete)
// 						 res.OnChainTime=time.Unix(queryResponse.Timestamp.Seconds,0).Format("2006-01-02 15:04:05")
// 						 results= append(results, res)
// 				 }
// 				 if err!=nil {
// 						 return nil,err
// 				 }
// 		 }
// 		 return results, nil
//  }
 
 
 func main() {
 
	 chaincode, err := contractapi.NewChaincode(new(SmartContract))
 
	 if err != nil {
		 fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		 return
	 }
 
	 if err := chaincode.Start(); err != nil {
		 fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	 }
 }
 
