import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance} from '@ant-design/pro-form';
import { ProFormTextArea } from '@ant-design/pro-form';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressActionStaticInfo,
  CtlIngressActionStaticSave,
  CtlIngressNamespaceMap,
  DataMapAppDataIngressActionStaticContentType,
  DataMapAppDataIngressActionStaticDataType,
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

  const xFormRef = useRef<ProFormInstance<API.AppDataIngressActionStatic>>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadNamespaceMapData = async () => {
    const dataList: any[] = [];

    const xReq = {};

    const xRespData = await CtlIngressNamespaceMap(xReq);

    if (xRespData && xRespData.data) {
      xRespData.data.forEach((xItem) => {
        dataList.push({ label: xItem.data_name, value: xItem.data_id });
      });
    }

    return dataList;
  };

  const loadDataInfo = async (id?: string) => {
    setIsLoading(true);

    if (id) {
      const respData = await CtlIngressActionStaticInfo({ data_id: id });
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

  const onSaveData = async (xData: API.AppDataIngressActionStatic) => {
    xData.data_id = props.id;

    const xResp = await CtlIngressActionStaticSave(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressActionStatic>
      title={'编辑资源'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressActionStatic);
      }}
      formRef={xFormRef}
    >
      <Spin spinning={xIsLoading} delay={100}>
        <ProFormSelect
          name="namespace_id"
          label="网关空间"
          placeholder="请选择网关空间"
          rules={[{ required: true, message: '请选择网关空间.' }]}
          request={loadNamespaceMapData}
          fieldProps={{ showSearch: true }}
        />

        <ProFormText
          name="title"
          label="资源名称"
          placeholder="请输入资源名称"
          rules={[{ required: true, message: '请输入资源名称.' }]}
        />

        <ProFormSelect
          name="content_type"
          label="资源类型"
          placeholder="请选择资源类型"
          fieldProps={{ options: DataMapAppDataIngressActionStaticContentType ,showSearch: true}}
          rules={[{ required: true, message: '请选择资源类型.' }]}
        />

        <ProFormSelect
          name="data_type"
          label="数据类型"
          valueEnum={DataMapAppDataIngressActionStaticDataType}
          placeholder="请选择数据类型"
          rules={[{ required: true, message: '请选择数据类型.' }]}
        />

        <ProFormTextArea
          name="data"
          label="数据内容"
          placeholder="请输入数据内容"
          rules={[{ required: true, message: '请输入数据内容.' }]}
        />
        <ProFormSelect
          name="state"
          label="状态"
          valueEnum={{
            enable: { text: '启用', status: 'enable' },
            disable: { text: '禁用', status: 'disable' },
          }}
          placeholder="请选择状态"
          rules={[{ required: true, message: '请选择状态.' }]}
        />
      </Spin>
    </ModalForm>
  );
};

export default EditForm;
