package com.herbertgao.telegram.database.service;

import com.herbertgao.telegram.database.entity.UserTemplate;
import com.herbertgao.telegram.database.entity.UserTemplateExample;
import com.herbertgao.telegram.database.mapper.UserTemplateMapper;
import com.herbertgao.telegram.util.IdUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @program: gaokao_bot
 * @description:
 * @author: HerbertGao
 * @create: 2020-04-08 18:13
 **/
@Service
public class UserTemplateService {

    @Autowired
    private UserTemplateMapper mapper;

    /**
     * 通过用户ID获取模板列表
     *
     * @param userId
     * @return
     */
    public List<UserTemplate> getUserTemplateListByUserId(Integer userId) {
        UserTemplateExample example = new UserTemplateExample();
        example.createCriteria()
                .andUserIdEqualTo(userId);
        return mapper.selectByExample(example);
    }

    /**
     * @param userId
     * @param examUsername
     */
    public Long insertDefaultTemplateWithUsername(Integer userId, String examUsername) {
        UserTemplate defaultTemplate = getDefaultTemplate();
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName("@" + examUsername);
        userTemplate.setTemplateContent(defaultTemplate.getTemplateContent() + "\n@" + examUsername);
        mapper.insert(userTemplate);
        return id;
    }

    /**
     * @param userId
     * @param template
     */
    public Long insertTemplateWithUsername(Integer userId, String template, String templateName) {
        UserTemplate userTemplate = new UserTemplate();
        Long id = IdUtil.getId();
        userTemplate.setId(id);
        userTemplate.setUserId(userId);
        userTemplate.setTemplateName(templateName);
        userTemplate.setTemplateContent(template);
        mapper.insert(userTemplate);
        return id;
    }

    /**
     * 获取默认模板列表
     *
     * @return
     */
    public UserTemplate getDefaultTemplate() {
        return getUserTemplateListByUserId(0).get(0);
    }

    /**
     * @param userId
     * @param id
     * @return
     */
    public UserTemplate getTemplateById(Integer userId, String id) {
        UserTemplateExample example = new UserTemplateExample();
        example.createCriteria()
                .andIdEqualTo(Long.parseLong(id))
                .andUserIdEqualTo(userId);
        List<UserTemplate> templateList = mapper.selectByExample(example);
        if (templateList.size() > 0) {
            return mapper.selectByExample(example).get(0);
        } else {
            return null;
        }
    }

    /**
     * @param userId
     * @param id
     */
    public void remove(Integer userId, String id) {
        UserTemplateExample example = new UserTemplateExample();
        example.createCriteria()
                .andIdEqualTo(Long.parseLong(id))
                .andUserIdEqualTo(userId);
        List<UserTemplate> templateList = mapper.selectByExample(example);
        if (templateList.size() > 0) {
            mapper.deleteByExample(example);
        }
    }

    /**
     * @param template
     */
    public void update(UserTemplate template) {
        mapper.updateByPrimaryKey(template);
    }
}
