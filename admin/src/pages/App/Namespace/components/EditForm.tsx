import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressNamespaceInfo,
  CtlIngressNamespaceSave,
  DataMapAppDataState,
} from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Spin } from 'antd';

export type EditFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  id?: string;
};

const EditForm: React.FC<EditFormProps> = (props) => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);

  const xFormRef = useRef<ProFormInstance<API.AppDataIngressNamespace>>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadDataInfo = async (id?: string) => {
    setIsLoading(true);

    if (id) {
      const respData = await CtlIngressNamespaceInfo({ data_id: id });
      if (respData && respData.error_code === '0') {
        if (respData.data) {
          const xDataInfo = respData.data;
          xFormRef.current?.setFieldsValue(xDataInfo);
        }
      }
    }

    setIsLoading(false);
  };

  // 加载
  useEffect(() => {
    if (props.modalVisible) {
      loadDataInfo(props.id);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.id, props.modalVisible]);

  const onSaveData = async (xData: API.AppDataIngressNamespace) => {
    xData.data_id = props.id;

    const xResp = await CtlIngressNamespaceSave(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressNamespace>
      title={'编辑集群空间'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressNamespace);
      }}
      formRef={xFormRef}
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProFormText
          name="data_id"
          label="空间ID"
          placeholder="请输入空间ID"
          disabled
          rules={[{ required: true, message: '请输入空间ID.' }]}
        />

        <ProFormText
          name="title"
          label="空间名称"
          placeholder="请输入空间名称"
          rules={[{ required: true, message: '请输入空间名称.' }]}
        />
        <ProFormSelect
          name="state"
          label="状态"
          valueEnum={DataMapAppDataState}
          placeholder="请选择状态"
          rules={[{ required: true, message: '请选择状态.' }]}
        />
      </Spin>
    </ModalForm>
  );
};

export default EditForm;
