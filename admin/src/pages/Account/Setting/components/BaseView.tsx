import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import ProForm, { ProFormText } from '@ant-design/pro-form';
import 'moment/locale/zh-cn';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Spin } from 'antd';
import { CtlUserGetCurrent, CtlUserUpdateInfo } from '@/services/ant-design-pro/api';

const BaseView: React.FC = () => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);

  const xFormRef = useRef<ProFormInstance<API.UserData>>();

  const loadDataInfo = async () => {
    setIsLoading(true);

    const respData = await CtlUserGetCurrent();
    if (respData && respData.error_code === '0') {
      if (respData.data) {
        const xDataInfo = respData.data;
        xFormRef.current?.setFieldsValue(xDataInfo);
      }
    }

    setIsLoading(false);
  };

  // 加载
  useEffect(() => {
    loadDataInfo();
  }, []);

  const onSaveData = async (xData: API.UserData) => {
    const xResp = await CtlUserUpdateInfo(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ProForm<API.UserData>
      onFinish={async (values) => {
        return onSaveData(values as API.UserData);
      }}
      submitter={{
        resetButtonProps: {
          style: {
            display: 'none',
          },
        },
      }}
      formRef={xFormRef}
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProFormText
          name="data_id"
          label="用户ID"
          placeholder="请输入用户ID"
          disabled
          rules={[{ required: true, message: '请输入用户ID.' }]}
        />

        <ProFormText
          name="username"
          label="用户名"
          placeholder="请输入用户名"
          disabled
          rules={[{ required: true, message: '请输入用户名.' }]}
        />

        <ProFormText
          name="nickname"
          label="用户昵称"
          placeholder="请输入用户昵称"
          rules={[{ required: true, message: '请输入用户昵称.' }]}
        />
        <ProFormText
          name="real_name"
          label="用户姓名"
          placeholder="请输入用户姓名"
          rules={[{ required: true, message: '请输入用户姓名.' }]}
        />
      </Spin>
    </ProForm>
  );
};

export default BaseView;
