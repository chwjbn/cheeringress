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
      title={'添加路由规则'}
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
        label="规则名称"
        placeholder="请输入规则名称"
        rules={[{ required: true, message: '请输入规则名称.' }]}
      />

      <ProFormDigit
        name="order_no"
        label="规则序号"
        placeholder="请输入规则序号"
        rules={[{ required: true, message: '请输入规则序号.' }]}
      />

      <ProFormSelect
        name="http_method"
        label="HTTP方法"
        valueEnum={DataMapHttpMethod}
        placeholder="请选择HTTP方法"
        rules={[{ required: true, message: '请选择HTTP方法.' }]}
      />

      <ProFormSelect
        name="match_target"
        label="规则匹配目标"
        valueEnum={DataMapHttpTargetItem}
        placeholder="请选择规则匹配目标"
        rules={[{ required: true, message: '请选择规则匹配目标.' }]}
      />

      <ProFormSelect
        name="match_op"
        label="规则匹配方式"
        valueEnum={DataMapRuleStringOp}
        placeholder="请选择规则匹配方式"
        rules={[{ required: true, message: '请选择规则匹配方式.' }]}
      />

      <ProFormText
        name="match_value"
        label="规则匹配内容"
        placeholder="请输入规则匹配内容"
        rules={[{ required: true, message: '请输入规则匹配内容.' }]}
      />

      <ProFormSelect
        name="action_type"
        label="响应类型"
        valueEnum={DataMapAppDataIngressActionType}
        placeholder="请选择响应类型"
        rules={[{ required: true, message: '请选择响应类型.' }]}
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
        label="响应内容"
        placeholder="请选择响应内容"
        fieldProps={{ options: xIngressActionValueMapData, showSearch: true }}
        rules={[{ required: true, message: '请选择响应内容.' }]}
      />
    </ModalForm>
  );
};

export default AddRuleForm;
