import React, { useEffect, useState } from 'react';
import { ProFormUploadDragger } from '@ant-design/pro-form';
import { CtlAliyunOSSUploadArg } from '@/services/ant-design-pro/api';
import type { UploadProps } from 'antd';
import type { UploadFile } from 'antd/es/upload/interface';

interface OSSDataType {
  dir?: string;
  host?: string;
  accessId?: string;
  policy?: string;
  signature?: string;
  baseDir?: string;
}

export type AliyunOSSUploadProps = {
  label?: string;
  width?: number | 'sm' | 'md' | 'xl' | 'xs' | 'lg' | undefined;
  fieldName?: string;
  fieldMessage?: string;
};

const AliyunOSSUpload: React.FC<AliyunOSSUploadProps> = (props) => {
  const [xOSSData, setOSSData] = useState<OSSDataType>();

  const [fileList, setFileList] = useState<UploadFile[]>();

  const initOSS = async () => {
    const xConfigResp = await CtlAliyunOSSUploadArg({});

    const xOSSConfig = {
      dir: 'devops/',
      host: xConfigResp.data?.endpoint,
      accessId: xConfigResp.data?.accessKeyId,
      policy: xConfigResp.data?.policy,
      signature: xConfigResp.data?.signature,
      baseDir: xConfigResp.data?.baseDir,
    };

    setOSSData(xOSSConfig);
  };

  useEffect(() => {
    initOSS();
  }, []);

  const getExtraData: UploadProps['data'] = (file) => ({
    key: file.url,
    OSSAccessKeyId: xOSSData?.accessId,
    policy: xOSSData?.policy,
    Signature: xOSSData?.signature,
  });

  const beforeUpload: UploadProps['beforeUpload'] = async (file) => {
    if (!xOSSData) {
      return false;
    }


    const filename = Date.now() + "_"+file.name;

    // @ts-ignore
    file.url = xOSSData.baseDir + '/' + filename;

    return file;
  };

  const handleChange: UploadProps['onChange'] = (info) => {
    let newFileList = [...info.fileList];

    newFileList = newFileList.map((file) => {


      if (file && (file.percent && file.percent >= 100)&&(file.status=='done')) {

        if(!file.url?.startsWith("http")){
            file.url = xOSSData?.host + '/' + file.url;
        }
 
      }

      return file;
    });

    //window.console.log(newFileList);

    setFileList(newFileList);
  };

  return (
    <ProFormUploadDragger
      label={props.label}
      width={props.width}
      name={props.fieldName}
      fieldProps={{
        beforeUpload: beforeUpload,
        data: getExtraData,
        name: 'file',
        onChange: handleChange,
        fileList: fileList,
      }}
      rules={[
        {
          required: true,
          message: props.fieldMessage,
        },
      ]}
      action={xOSSData?.host}
    />
  );
};

export default AliyunOSSUpload;
