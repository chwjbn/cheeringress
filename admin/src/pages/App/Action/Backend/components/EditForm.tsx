import React, { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import { CtlIngressActionBackendInfo, CtlIngressActionBackendSave, CtlIngressNamespaceMap, DataMapAppDataIngressActionBackendBalanceType } from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';
import { Spin } from 'antd';

export type EditFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  id?: string;
};

const EditForm: React.FC<EditFormProps> = (props) => {
  const [xIsLoading, setIsLoading] = useState<boolean>(true);

  const xFormRef = useRef<ProFormInstance<API.AppDataIngressActionBackend>>();

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
      const respData = await CtlIngressActionBackendInfo({ data_id: id });
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

  const onSaveData = async (xData: API.AppDataIngressActionBackend) => {
    xData.data_id = props.id;

    const xResp = await CtlIngressActionBackendSave(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressActionBackend>
      title={'编辑反向代理'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressActionBackend);
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
        label="代理名称"
        placeholder="请输入代理名称"
        rules={[{ required: true, message: '请输入代理名称.' }]}
      />

      <ProFormSelect
        name="balance_type"
        label="负载均衡策略"
        valueEnum={DataMapAppDataIngressActionBackendBalanceType}
        placeholder="请选择负载均衡策略"
        rules={[{ required: true, message: '请选择负载均衡策略.' }]}
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


