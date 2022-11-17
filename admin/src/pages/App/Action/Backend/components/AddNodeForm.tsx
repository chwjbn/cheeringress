import React from 'react';
import { ModalForm, ProFormDigit, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import { CtlIngressActionBackendNodeAdd } from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';

export type AddNodeFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  namespace_id?: string;
  backend_id?: string;
};

const AddNodeForm: React.FC<AddNodeFormProps> = (props) => {
  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const onSaveData = async (xData: API.AppDataIngressActionBackendNode) => {

    xData.namespace_id=props.namespace_id;
    xData.backend_id=props.backend_id;

    const xResp = await CtlIngressActionBackendNodeAdd(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressActionBackendNode>
      title={'添加代理节点'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressActionBackendNode);
      }}
    >
      <ProFormText
        name="title"
        label="节点名称"
        placeholder="请输入节点名称"
        rules={[{ required: true, message: '请输入节点名称.' }]}
      />

      <ProFormText
        name="server_host"
        label="服务器地址"
        placeholder="请输入服务器地址"
        rules={[{ required: true, message: '请输入服务器地址.' }]}
      />

      <ProFormDigit
        name="server_port"
        label="服务器端口"
        placeholder="请输入服务器端口"
        min={0}
        max={65535}
        rules={[{ required: true, message: '请输入服务器端口.' }]}
        fieldProps={{ precision: 0 }}
      />

      <ProFormDigit
        name="weight_score"
        label="服务器权重"
        placeholder="请输入服务器权重"
        min={1}
        max={100}
        rules={[{ required: true, message: '请输入服务器权重.' }]}
        fieldProps={{ precision: 0 }}
      />
    </ModalForm>
  );
};

export default AddNodeForm;
