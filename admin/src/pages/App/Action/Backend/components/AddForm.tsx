import React from 'react';
import { ModalForm, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressActionBackendAdd,
  CtlIngressNamespaceMap,
  DataMapAppDataIngressActionBackendBalanceType,
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

  const onSaveData = async (xData: API.AppDataIngressActionBackend) => {
    const xResp = await CtlIngressActionBackendAdd(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressActionBackend>
      title={'添加反向代理'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressActionBackend);
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
    </ModalForm>
  );
};

export default AddForm;
