// Code generated by counterfeiter. DO NOT EDIT.
package internalfakes

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vmware-labs/marketplace-cli/v2/internal"
)

type FakeS3Client struct {
	PutObjectStub        func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	putObjectMutex       sync.RWMutex
	putObjectArgsForCall []struct {
		arg1 context.Context
		arg2 *s3.PutObjectInput
		arg3 []func(*s3.Options)
	}
	putObjectReturns struct {
		result1 *s3.PutObjectOutput
		result2 error
	}
	putObjectReturnsOnCall map[int]struct {
		result1 *s3.PutObjectOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeS3Client) PutObject(arg1 context.Context, arg2 *s3.PutObjectInput, arg3 ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	fake.putObjectMutex.Lock()
	ret, specificReturn := fake.putObjectReturnsOnCall[len(fake.putObjectArgsForCall)]
	fake.putObjectArgsForCall = append(fake.putObjectArgsForCall, struct {
		arg1 context.Context
		arg2 *s3.PutObjectInput
		arg3 []func(*s3.Options)
	}{arg1, arg2, arg3})
	fake.recordInvocation("PutObject", []interface{}{arg1, arg2, arg3})
	fake.putObjectMutex.Unlock()
	if fake.PutObjectStub != nil {
		return fake.PutObjectStub(arg1, arg2, arg3...)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.putObjectReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeS3Client) PutObjectCallCount() int {
	fake.putObjectMutex.RLock()
	defer fake.putObjectMutex.RUnlock()
	return len(fake.putObjectArgsForCall)
}

func (fake *FakeS3Client) PutObjectCalls(stub func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)) {
	fake.putObjectMutex.Lock()
	defer fake.putObjectMutex.Unlock()
	fake.PutObjectStub = stub
}

func (fake *FakeS3Client) PutObjectArgsForCall(i int) (context.Context, *s3.PutObjectInput, []func(*s3.Options)) {
	fake.putObjectMutex.RLock()
	defer fake.putObjectMutex.RUnlock()
	argsForCall := fake.putObjectArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeS3Client) PutObjectReturns(result1 *s3.PutObjectOutput, result2 error) {
	fake.putObjectMutex.Lock()
	defer fake.putObjectMutex.Unlock()
	fake.PutObjectStub = nil
	fake.putObjectReturns = struct {
		result1 *s3.PutObjectOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeS3Client) PutObjectReturnsOnCall(i int, result1 *s3.PutObjectOutput, result2 error) {
	fake.putObjectMutex.Lock()
	defer fake.putObjectMutex.Unlock()
	fake.PutObjectStub = nil
	if fake.putObjectReturnsOnCall == nil {
		fake.putObjectReturnsOnCall = make(map[int]struct {
			result1 *s3.PutObjectOutput
			result2 error
		})
	}
	fake.putObjectReturnsOnCall[i] = struct {
		result1 *s3.PutObjectOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeS3Client) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.putObjectMutex.RLock()
	defer fake.putObjectMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeS3Client) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ internal.S3Client = new(FakeS3Client)