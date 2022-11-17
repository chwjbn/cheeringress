import React from 'react';
import { ModalForm, ProFormSelect, ProFormText, ProFormTextArea } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressActionStaticAdd,
  CtlIngressNamespaceMap,
  DataMapAppDataIngressActionStaticContentType,
  DataMapAppDataIngressActionStaticDataType,
} from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';

export type AddFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
};

const AddForm: React.FC<AddFormProps> = (props) => {
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

  const onSaveData = async (xData: API.AppDataIngressActionStatic) => {
    const xResp = await CtlIngressActionStaticAdd(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressActionStatic>
      title={'添加资源'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressActionStatic);
      }}
    >
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
        fieldProps={{ options: DataMapAppDataIngressActionStaticContentType,showSearch: true }}
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
    </ModalForm>
  );
};

export default AddForm;
