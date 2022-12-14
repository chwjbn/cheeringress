import React, { useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-form';
import { ModalForm, ProFormDigit, ProFormSelect, ProFormText } from '@ant-design/pro-form';
import moment from 'moment';
import 'moment/locale/zh-cn';
import {
  CtlIngressActionBackendMap,
  CtlIngressActionStaticMap,
  CtlIngressSiteRuleAdd,
  DataMapAppDataIngressActionType,
  DataMapHttpMethod,
  DataMapHttpTargetItem,
  DataMapRuleStringOp,
} from '@/services/ant-design-pro/api';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';

export type AddRuleFormProps = {
  onModalVisible: (flag: boolean) => void;
  modalVisible: boolean;
  namespace_id?: string;
  site_id?: string;
};

const AddRuleForm: React.FC<AddRuleFormProps> = (props) => {
  const xFormRef = useRef<ProFormInstance<API.AppDataIngressSite>>();
  const [xIngressActionValueMapData, setIngressActionValueMapData] = useState<any>();

  const onFormVisible = (flag: boolean) => {
    moment.locale('zh-cn');
    return props.onModalVisible(flag);
  };

  const loadIngressActionValueMapData = async () => {
    const xDataList: any[] = [];

    const xActionType = xFormRef.current?.getFieldValue('action_type');

    const xNamespaceId = props.namespace_id;
    const xReq = { data_id: xNamespaceId };

    if (xActionType == 'backend') {
      const xRespData = await CtlIngressActionBackendMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    if (xActionType == 'static') {
      const xRespData = await CtlIngressActionStaticMap(xReq);
      if (xRespData && xRespData.data) {
        xRespData.data.forEach((xItem) => {
          xDataList.push({ label: xItem.data_name, value: xItem.data_id });
        });
      }
    }

    setIngressActionValueMapData(xDataList);
  };

  const onSaveData = async (xData: API.AppDataIngressSiteRule) => {
    xData.namespace_id = props.namespace_id;
    xData.site_id = props.site_id;

    const xResp = await CtlIngressSiteRuleAdd(xData);

    if (xResp.error_code === '0') {
      ShowSuccessMessage(xResp.error_message);
      return true;
    }

    ShowErrorMessage(xResp.error_message);

    return false;
  };

  return (
    <ModalForm<API.AppDataIngressSiteRule>
      title={'??????????????????'}
      width="80%"
      onVisibleChange={onFormVisible}
      visible={props.modalVisible}
      modalProps={{ destroyOnClose: true }}
      onFinish={async (values) => {
        return onSaveData(values as API.AppDataIngressSiteRule);
      }}
      formRef={xFormRef}
    >
      <ProFormText
        name="title"
        label="????????????"
        placeholder="?????????????????????"
        rules={[{ required: true, message: '?????????????????????.' }]}
      />

      <ProFormDigit
        name="order_no"
        label="????????????"
        placeholder="?????????????????????"
        rules={[{ required: true, message: '?????????????????????.' }]}
      />

      <ProFormSelect
        name="http_method"
        label="HTTP??????"
        valueEnum={DataMapHttpMethod}
        placeholder="?????????HTTP??????"
        rules={[{ required: true, message: '?????????HTTP??????.' }]}
      />

      <ProFormSelect
        name="match_target"
        label="??????????????????"
        valueEnum={DataMapHttpTargetItem}
        placeholder="???????????????????????????"
        rules={[{ required: true, message: '???????????????????????????.' }]}
      />

      <ProFormSelect
        name="match_op"
        label="??????????????????"
        valueEnum={DataMapRuleStringOp}
        placeholder="???????????????????????????"
        rules={[{ required: true, message: '???????????????????????????.' }]}
      />

      <ProFormText
        name="match_value"
        label="??????????????????"
        placeholder="???????????????????????????"
        rules={[{ required: true, message: '???????????????????????????.' }]}
      />

      <ProFormSelect
        name="action_type"
        label="????????????"
        valueEnum={DataMapAppDataIngressActionType}
        placeholder="?????????????????????"
        rules={[{ required: true, message: '?????????????????????.' }]}
        fieldProps={{
          onChange() {
            const xFormDataVal = xFormRef.current?.getFieldsValue();

            if (xFormDataVal) {
              xFormDataVal.action_value = undefined;
              xFormRef.current?.setFieldsValue(xFormDataVal);
            }

            loadIngressActionValueMapData();
          },
        }}
      />

      <ProFormSelect
        name="action_value"
        label="????????????"
        placeholder="?????????????????????"
        fieldProps={{ options: xIngressActionValueMapData, showSearch: true }}
        rules={[{ required: true, message: '?????????????????????.' }]}
      />
    </ModalForm>
  );
};

export default AddRuleForm;
