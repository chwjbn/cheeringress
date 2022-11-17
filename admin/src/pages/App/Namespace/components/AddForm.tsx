import React from 'react';
import { ModalForm, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import { CtlIngressNamespaceAdd } from '@/services/ant-design-pro/api';
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

  const onSaveData = async (xData: API.AppDataIngressNamespace) => {
    const xResp = await CtlIngressNamespaceAdd(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressNamespace>
      title={'添加集群空间'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressNamespace);
      }}
    >
      <ProFormText
        name="title"
        label="空间名称"
        placeholder="请输入空间名称"
        rules={[{ required: true, message: '请输入空间名称.' }]}
      />


    </ModalForm>
  );
};

export default AddForm;
