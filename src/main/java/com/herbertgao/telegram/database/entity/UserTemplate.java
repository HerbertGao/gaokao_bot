package com.herbertgao.telegram.database.entity;

public class UserTemplate extends UserTemplateKey {
    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column user_template.user_id
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    private Integer userId;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column user_template.template_name
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    private String templateName;

    /**
     *
     * This field was generated by MyBatis Generator.
     * This field corresponds to the database column user_template.template_content
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    private String templateContent;

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column user_template.user_id
     *
     * @return the value of user_template.user_id
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public Integer getUserId() {
        return userId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column user_template.user_id
     *
     * @param userId the value for user_template.user_id
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public void setUserId(Integer userId) {
        this.userId = userId;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column user_template.template_name
     *
     * @return the value of user_template.template_name
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public String getTemplateName() {
        return templateName;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column user_template.template_name
     *
     * @param templateName the value for user_template.template_name
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public void setTemplateName(String templateName) {
        this.templateName = templateName == null ? null : templateName.trim();
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method returns the value of the database column user_template.template_content
     *
     * @return the value of user_template.template_content
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public String getTemplateContent() {
        return templateContent;
    }

    /**
     * This method was generated by MyBatis Generator.
     * This method sets the value of the database column user_template.template_content
     *
     * @param templateContent the value for user_template.template_content
     *
     * @mbg.generated Wed Apr 08 18:12:30 CST 2020
     */
    public void setTemplateContent(String templateContent) {
        this.templateContent = templateContent == null ? null : templateContent.trim();
    }
}